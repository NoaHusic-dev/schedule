package models

type WeatherResponse struct {
	Temperature int `json:"temperature"`
	WeatherCode int `json:"weatherCode"`
}
