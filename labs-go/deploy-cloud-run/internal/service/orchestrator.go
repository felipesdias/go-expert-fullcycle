package service

import (
	"errors"
	"math"
	"regexp"

	"weather-api/internal/entity"
)

var (
	ErrInvalidZipCode  = errors.New("invalid zipcode")
	ErrZipCodeNotFound = errors.New("can not find zipcode")
)

type Orchestrator interface {
	GetWeather(cep string) (*entity.WeatherOutput, error)
}

type OrchestratorService struct {
	LocationProvider LocationProvider
	WeatherProvider  WeatherProvider
}

func NewOrchestratorService(lp LocationProvider, wp WeatherProvider) *OrchestratorService {
	return &OrchestratorService{
		LocationProvider: lp,
		WeatherProvider:  wp,
	}
}

func (s *OrchestratorService) GetWeather(cep string) (*entity.WeatherOutput, error) {
	if !isValidCEP(cep) {
		return nil, ErrInvalidZipCode
	}

	location, err := s.LocationProvider.GetLocation(cep)
	if err != nil {
		return nil, err
	}

	if location == nil || location.Localidade == "" || location.IsErro() {
		return nil, ErrZipCodeNotFound
	}

	weather, err := s.WeatherProvider.GetWeather(location.Localidade)
	if err != nil {
		return nil, err
	}

	tempC := weather.Current.TempC
	tempF := math.Round((tempC*1.8+32)*100) / 100
	tempK := math.Round((tempC+273)*100) / 100

	return &entity.WeatherOutput{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}
