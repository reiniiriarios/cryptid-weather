package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type WeatherData struct {
	TempC      float32
	FeelsLikeC float32
	Humidity   uint8
	Condition  string
	Code       uint16
}

func getCurrentWeather() (*WeatherData, error) {
	resp, err := getWeather("current")
	if err != nil {
		return nil, err
	}
	// Adapt
	return &WeatherData{
		TempC:      resp.Current.TempC,
		FeelsLikeC: resp.Current.FeelsLikeC,
		Humidity:   resp.Current.Humidity,
		Condition:  getStringFromWeatherCode(resp.Current.Condition.Code),
		Code:       resp.Current.Condition.Code,
	}, nil
}

func getWeather(endpoint string) (*weatherResponse, error) {
	// Fetch
	latlng := os.Getenv("WEATHER_API_LAT") + "," + os.Getenv("WEATHER_API_LNG")
	query := "?key=" + os.Getenv("WEATHER_API_KEY") + "&q=" + latlng
	resp, err := http.Get("http://api.weatherapi.com/v1/" + endpoint + ".json" + query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	// Read
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Parse
	var data weatherResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
