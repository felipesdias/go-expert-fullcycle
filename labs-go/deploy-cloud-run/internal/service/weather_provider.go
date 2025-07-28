package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"weather-api/internal/entity"
)

type WeatherProvider interface {
	GetWeather(city string) (*entity.WeatherApiOutput, error)
}

type WeatherApiProvider struct {
	BaseURL string
	ApiKey  string
}

func NewWeatherApiProvider(apiKey string) *WeatherApiProvider {
	return &WeatherApiProvider{
		BaseURL: "https://api.weatherapi.com/v1",
		ApiKey:  apiKey,
	}
}

func (p *WeatherApiProvider) GetWeather(city string) (*entity.WeatherApiOutput, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s/current.json?key=%s&q=%s", p.BaseURL, p.ApiKey, encodedCity)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather api request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data entity.WeatherApiOutput
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
