package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type WeatherData struct {
	TempC float32
}

func getCurrentWeather() (*WeatherData, error) {
	return getWeather("current")
}

func getWeather(endpoint string) (*WeatherData, error) {
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
	println(string(body))
	// Parse
	var data weatherResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	// Adapt
	return &WeatherData{
		TempC: data.Current.TempC,
	}, nil
}
