package service

import (
	"errors"
	"testing"

	"weather-api/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLocationProvider struct {
	mock.Mock
}

func (m *MockLocationProvider) GetLocation(cep string) (*entity.ViaCepOutput, error) {
	args := m.Called(cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ViaCepOutput), args.Error(1)
}

type MockWeatherProvider struct {
	mock.Mock
}

func (m *MockWeatherProvider) GetWeather(city string) (*entity.WeatherApiOutput, error) {
	args := m.Called(city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WeatherApiOutput), args.Error(1)
}

func TestOrchestratorService_GetWeather(t *testing.T) {
	t.Run("should return weather on success", func(t *testing.T) {
		mockLP := new(MockLocationProvider)
		mockWP := new(MockWeatherProvider)
		service := NewOrchestratorService(mockLP, mockWP)

		cep := "36570000"
		city := "Vicosa"
		tempC := 25.0

		mockLP.On("GetLocation", cep).Return(&entity.ViaCepOutput{Localidade: city}, nil)
		mockWP.On("GetWeather", city).Return(&entity.WeatherApiOutput{Current: entity.CurrentWeather{TempC: tempC}}, nil)

		expected := &entity.WeatherOutput{
			TempC: 25.0,
			TempF: 77.0,
			TempK: 298.0,
		}

		result, err := service.GetWeather(cep)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockLP.AssertExpectations(t)
		mockWP.AssertExpectations(t)
	})

	t.Run("should return error for invalid cep", func(t *testing.T) {
		service := NewOrchestratorService(nil, nil)
		_, err := service.GetWeather("12345")
		assert.EqualError(t, err, "invalid zipcode")
	})

	t.Run("should return error when cep not found", func(t *testing.T) {
		mockLP := new(MockLocationProvider)
		service := NewOrchestratorService(mockLP, nil)
		cep := "00000000"

		mockLP.On("GetLocation", cep).Return(nil, nil)

		_, err := service.GetWeather(cep)
		assert.EqualError(t, err, "can not find zipcode")
	})

	t.Run("should return error when location provider fails", func(t *testing.T) {
		mockLP := new(MockLocationProvider)
		service := NewOrchestratorService(mockLP, nil)
		cep := "36570000"
		expectedErr := errors.New("provider error")

		mockLP.On("GetLocation", cep).Return(nil, expectedErr)

		_, err := service.GetWeather(cep)
		assert.EqualError(t, err, expectedErr.Error())
	})
}
