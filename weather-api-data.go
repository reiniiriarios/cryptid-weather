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
	Condition        struct {
		Text string `json:"text"` // current conditions
		Icon string `json:"icon"` // url
		Code uint16 `json:"code"` // condition unique code
	} `json:"condition"`
	WindMph    float32 `json:"wind_mph"`
	WindKph    float32 `json:"wind_kph"`
	WindDegree uint16  `json:"wind_degree"`
	WindDir    string  `json:"wind_dir"`
	PressureMb float32 `json:"pressure_mb"`
	PressureIn float32 `json:"pressure_in"`
	Humidity   uint8   `json:"humidity"` // percentage
	Cloud      uint8   `json:"cloud"`    // percentage
	IsDay      uint8   `json:"is_day"`   // 1 = yes, 0 = no
	UV         float32 `json:"uv"`
	GustMph    float32 `json:"gust_mph"`
	GustKph    float32 `json:"gust_kph"`
}

// https://www.weatherapi.com/docs/weather_conditions.json
var weatherCodes = map[uint16]string{
	1000: "clear",
	1003: "partly_cloudy",
	1006: "cloudy",
	1009: "overcast",
	1030: "mist",
	1063: "patchy_rain_possible",
	1066: "patchy_snow_possible",
	1069: "patchy_sleet_possible",
	1072: "patchy_freezing_drizzle_possible",
	1087: "thundery_outbreaks_possible",
	1114: "blowing_snow",
	1117: "blizzard",
	1135: "fog",
	1147: "freezing_fog",
	1150: "patchy_light_drizzle",
	1153: "light_drizzle",
	1168: "freezing_drizzle",
	1171: "heavy_freezing_drizzle",
	1180: "patchy_light_rain",
	1183: "light_rain",
	1186: "moderate_rain_at_times",
	1189: "moderate_rain",
	1192: "heavy_rain_at_times",
	1195: "heavy_rain",
	1198: "light_freezing_rain",
	1201: "moderate_or_heavy_freezing_rain",
	1204: "light_sleet",
	1207: "moderate_or_heavy_sleet",
	1210: "patchy_light_snow",
	1213: "light_snow",
	1216: "patchy_moderate_snow",
	1219: "moderate_snow",
	1222: "patchy_heavy_snow",
	1225: "heavy_snow",
	1237: "ice_pellets",
	1240: "light_rain_shower",
	1243: "moderate_or_heavy_rain_shower",
	1246: "torrential_rain_shower",
	1249: "light_sleet_showers",
	1252: "moderate_or_heavy_sleet_showers",
	1255: "light_snow_showers",
	1258: "moderate_or_heavy_snow_showers",
	1261: "light_showers_of_ice_pellets",
	1264: "moderate_or_heavy_showers_of_ice_pellets",
	1273: "patchy_light_rain_with_thunder",
	1276: "moderate_or_heavy_rain_with_thunder",
	1279: "patchy_light_snow_with_thunder",
	1282: "moderate_or_heavy_snow_with_thunder",
}

func getStringFromWeatherCode(code uint16) string {
	if s, exists := weatherCodes[code]; exists {
		return s
	}
	return "unknown"
}
