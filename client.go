package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	Namespace               = "feed.weatherapi"
	currentForecastEndpoint = "/current.json"
)

// BindEnvs registers environment variable mappings for this feed's config namespace.
func BindEnvs(v *viper.Viper) {
	_ = v.BindEnv("feed.weatherapi.url", "FEED_WEATHERAPI_URL")
	_ = v.BindEnv("feed.weatherapi.key", "FEED_WEATHERAPI_KEY")
}

// Client for the WeatherAPI feed
type Client struct {
	logger *zap.Logger
	URL    string
	Key    string
}

// NewClient creates a new WeatherAPI client
func NewClient(logger *zap.Logger, cfg *Config) (*Client, error) {
	return &Client{
		logger: logger.Named(Namespace),
		URL:    cfg.URL,
		Key:    cfg.Key,
	}, nil
}

// GetCurrentForecast retrieves the current weather forecast based on latitude and longitude.
func (c *Client) GetCurrentForecast(latitude float64, longitude float64) (*Payload, error) {
	req, err := c.getRequest(currentForecastEndpoint)
	if err != nil {
		return nil, err
	}

	// Add query parameters to the request.
	c.addQueryParams(req, map[string]string{
		"q":   fmt.Sprintf("%.8f,%.8f", latitude, longitude),
		"key": c.Key,
	})

	// Send the request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	// Read the response body.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response.
	var output Payload
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &output, nil
}

// getRequest creates a new HTTP GET request for a specific API path.
func (c *Client) getRequest(path string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.URL, path)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return req, nil
}

// addQueryParams adds query parameters to the given request.
func (c *Client) addQueryParams(req *http.Request, params map[string]string) {
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}

	req.URL.RawQuery = query.Encode()
}
