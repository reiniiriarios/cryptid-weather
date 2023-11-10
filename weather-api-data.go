package main

type weatherResponse struct {
	Location weatherLocation        `json:"location"`
	Current  currentWeatherResponse `json:"current"`
}

type weatherLocation struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
	Timezone       string  `json:"tz_id"`
	LocalTimeEpoch float64 `json:"localtime_epoch"`
	LocalTime      string  `json:"localtime"`
}

type currentWeatherResponse struct {
	LastUpdated      string  `json:"last_updated"`
	LastUpdatedEpoch int     `json:"last_updated_epoch"`
	TempC            float32 `json:"temp_c"`
	TempF            float32 `json:"temp_f"`
	FeelsLikeC       float32 `json:"feelslike_c"`
	FeelsLikeF       float32 `json:"feelslike_f"`
	Text             string  `json:"current:text"` // current conditions
	Icon             string  `json:"current:icon"` // url
	Code             string  `json:"current:code"` // condition unique code
	WindMph          float32 `json:"wind_mph"`
	WindKph          float32 `json:"wind_kph"`
	WindDegree       int     `json:"wind_degree"`
	WindDir          string  `json:"wind_dir"`
	PressureMb       float32 `json:"pressure_mb"`
	PressureIn       float32 `json:"pressure_in"`
	Humidity         int     `json:"humidity"` // percentage
	Cloud            int     `json:"cloud"`    // percentage
	IsDay            int     `json:"is_day"`   // 1 = yes, 0 = no
	UV               float32 `json:"uv"`
	GustMph          float32 `json:"gust_mph"`
	GustKph          float32 `json:"gust_kph"`
}
