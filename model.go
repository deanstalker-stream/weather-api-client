package weatherapi

type Config struct {
	URL string `json:"url" mapstructure:"url"`
	Key string `json:"key" mapstructure:"key"`
}

// Payload represents the response from the weatherapi.com API
type Payload struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

// Location represents the location of a weather observation
type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TimezoneID     string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

// Current represents the current weather conditions
type Current struct {
	LastUpdatedEpoch    int64     `json:"last_updated_epoch"`
	LastUpdated         string    `json:"last_updated"`
	TempCelsius         float64   `json:"temp_c"`
	TempFahrenheit      float64   `json:"temp_f"`
	IsDay               int64     `json:"is_day"`
	Condition           Condition `json:"condition"`
	WindMPH             float64   `json:"wind_mph"`
	WindKPH             float64   `json:"wind_kph"`
	WindDegree          float64   `json:"wind_degree"`
	WindDirection       string    `json:"wind_dir"`
	PressureMB          float64   `json:"pressure_mb"`
	PressureIN          float64   `json:"pressure_in"`
	PrecipMM            float64   `json:"precip_mm"`
	PrecipIN            float64   `json:"precip_in"`
	Humidity            float64   `json:"humidity"`
	Cloud               float64   `json:"cloud"`
	FeelsLikeCelsius    float64   `json:"feelslike_c"`
	FeelsLikeFahrenheit float64   `json:"feelslike_f"`
	WindchillCelsius    float64   `json:"windchill_c"`
	WindchillFahrenheit float64   `json:"windchill_f"`
	HeatindexCelsius    float64   `json:"heatindex_c"`
	HeatindexFahrenheit float64   `json:"heatindex_f"`
	DewpointCelsius     float64   `json:"dewpoint_c"`
	DewpointFahrenheit  float64   `json:"dewpoint_f"`
	VisibilityKM        float64   `json:"vis_km"`
	VisibilityMiles     float64   `json:"vis_miles"`
	UV                  float64   `json:"uv"`
	GustMPH             float64   `json:"gust_mph"`
	GustKPH             float64   `json:"gust_kph"`
}

// Condition represents the weather condition
type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int64  `json:"code"`
}
