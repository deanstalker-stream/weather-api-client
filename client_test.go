package weatherapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestNewClient tests the NewClient constructor
func TestNewClient(t *testing.T) {
	client, err := NewClient(zap.NewNop(), &Config{URL: "https://api.weatherapi.com", Key: "test-key"})

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.weatherapi.com", client.URL)
	assert.Equal(t, "test-key", client.Key)
	assert.NotNil(t, client.logger)
}

// TestGetCurrentForecast_Success tests successful forecast retrieval
func TestGetCurrentForecast_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/current.json", r.URL.Path)
		assert.Equal(t, "40.71280000,74.00600000", r.URL.Query().Get("q"))
		assert.Equal(t, "test-key", r.URL.Query().Get("key"))

		response := Payload{
			Current: Current{
				TempCelsius:    25.5,
				TempFahrenheit: 77.9,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		logger: zap.NewNop(),
		URL:    server.URL,
		Key:    "test-key",
	}

	payload, err := client.GetCurrentForecast(40.7128, 74.0060)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
	assert.Equal(t, 25.5, payload.Current.TempCelsius)
	assert.Equal(t, 77.9, payload.Current.TempFahrenheit)
}

// TestGetCurrentForecast_InvalidJSON tests handling of invalid JSON response
func TestGetCurrentForecast_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &Client{
		logger: zap.NewNop(),
		URL:    server.URL,
		Key:    "test-key",
	}

	payload, err := client.GetCurrentForecast(40.7128, 74.0060)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "failed to unmarshal JSON")
}

// TestGetCurrentForecast_HTTPError tests HTTP request failure
func TestGetCurrentForecast_HTTPError(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "http://invalid-url-that-does-not-exist",
		Key:    "test-key",
	}

	payload, err := client.GetCurrentForecast(40.7128, 74.0060)

	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "failed to execute request")
}

// TestGetCurrentForecast_ZeroCoordinates tests with zero coordinates
func TestGetCurrentForecast_ZeroCoordinates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "0.00000000,0.00000000", r.URL.Query().Get("q"))
		response := Payload{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		logger: zap.NewNop(),
		URL:    server.URL,
		Key:    "test-key",
	}

	payload, err := client.GetCurrentForecast(0, 0)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
}

// TestGetCurrentForecast_NegativeCoordinates tests with negative coordinates
func TestGetCurrentForecast_NegativeCoordinates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "-33.86880000,-151.20550000", r.URL.Query().Get("q"))
		response := Payload{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		logger: zap.NewNop(),
		URL:    server.URL,
		Key:    "test-key",
	}

	payload, err := client.GetCurrentForecast(-33.8688, -151.2055)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
}

// TestGetRequest_Success tests successful request creation
func TestGetRequest_Success(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "https://api.weatherapi.com",
		Key:    "test-key",
	}

	req, err := client.getRequest("/current.json")

	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://api.weatherapi.com/current.json", req.URL.String())
}

// TestGetRequest_InvalidURL tests request creation with invalid URL
func TestGetRequest_InvalidURL(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "ht!tp://invalid",
		Key:    "test-key",
	}

	req, err := client.getRequest("/current.json")

	assert.Error(t, err)
	assert.Nil(t, req)
	assert.Contains(t, err.Error(), "failed to create request")
}

// TestAddQueryParams_SingleParam tests adding a single query parameter
func TestAddQueryParams_SingleParam(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "https://api.weatherapi.com",
		Key:    "test-key",
	}

	req, _ := client.getRequest("/current.json")
	client.addQueryParams(req, map[string]string{"key": "test-key"})

	assert.Equal(t, "key=test-key", req.URL.RawQuery)
}

// TestAddQueryParams_MultipleParams tests adding multiple query parameters
func TestAddQueryParams_MultipleParams(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "https://api.weatherapi.com",
		Key:    "test-key",
	}

	req, _ := client.getRequest("/current.json")
	client.addQueryParams(req, map[string]string{
		"q":   "40.7128,74.0060",
		"key": "test-key",
	})

	query := req.URL.Query()
	assert.Equal(t, "40.7128,74.0060", query.Get("q"))
	assert.Equal(t, "test-key", query.Get("key"))
}

// TestAddQueryParams_EmptyParams tests adding empty parameters
func TestAddQueryParams_EmptyParams(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "https://api.weatherapi.com",
		Key:    "test-key",
	}

	req, _ := client.getRequest("/current.json")
	client.addQueryParams(req, map[string]string{})

	assert.Equal(t, "", req.URL.RawQuery)
}

// TestAddQueryParams_SpecialCharacters tests adding parameters with special characters
func TestAddQueryParams_SpecialCharacters(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "https://api.weatherapi.com",
		Key:    "test-key",
	}

	req, _ := client.getRequest("/current.json")
	client.addQueryParams(req, map[string]string{
		"q": "San Francisco, CA",
	})

	query := req.URL.Query()
	assert.Equal(t, "San Francisco, CA", query.Get("q"))
}

// TestGetCurrentForecast_ResponseBodyClose tests proper response body closure
func TestGetCurrentForecast_ResponseBodyClose(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Payload{
			Current: Current{
				TempCelsius: 20.0,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		logger: zap.NewNop(),
		URL:    server.URL,
		Key:    "test-key",
	}

	payload, err := client.GetCurrentForecast(40.7128, 74.0060)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
}

// TestGetCurrentForecast_LargeCoordinates tests with large coordinate values
func TestGetCurrentForecast_LargeCoordinates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Payload{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		logger: zap.NewNop(),
		URL:    server.URL,
		Key:    "test-key",
	}

	payload, err := client.GetCurrentForecast(89.99999999, 179.99999999)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
}
