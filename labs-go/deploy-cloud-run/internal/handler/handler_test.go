package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-api/internal/entity"
	"weather-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrchestratorService struct {
	mock.Mock
}

func (m *MockOrchestratorService) GetWeather(cep string) (*entity.WeatherOutput, error) {
	args := m.Called(cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WeatherOutput), args.Error(1)
}

func TestWeatherHandler_GetWeather(t *testing.T) {
	setup := func(cep string, mockSvc *MockOrchestratorService) *httptest.ResponseRecorder {
		handler := NewWeatherHandler(mockSvc)
		req := httptest.NewRequest("GET", "/weather/"+cep, nil)
		rr := httptest.NewRecorder()

		router := chi.NewRouter()
		router.Get("/weather/{cep}", handler.GetWeather)
		router.ServeHTTP(rr, req)
		return rr
	}

	t.Run("should return 200 OK with weather data", func(t *testing.T) {
		mockSvc := new(MockOrchestratorService)
		cep := "36570000"
		weather := &entity.WeatherOutput{TempC: 20, TempF: 68, TempK: 293}

		mockSvc.On("GetWeather", cep).Return(weather, nil)
		rr := setup(cep, mockSvc)

		assert.Equal(t, http.StatusOK, rr.Code)

		var body entity.WeatherOutput
		err := json.Unmarshal(rr.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Equal(t, *weather, body)
		mockSvc.AssertExpectations(t)
	})

	t.Run("should return 422 for invalid zipcode", func(t *testing.T) {
		mockSvc := new(MockOrchestratorService)
		cep := "123"

		mockSvc.On("GetWeather", cep).Return(nil, service.ErrInvalidZipCode)
		rr := setup(cep, mockSvc)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Equal(t, "invalid zipcode\n", rr.Body.String())
		mockSvc.AssertExpectations(t)
	})

	t.Run("should return 404 for zipcode not found", func(t *testing.T) {
		mockSvc := new(MockOrchestratorService)
		cep := "00000000"

		mockSvc.On("GetWeather", cep).Return(nil, service.ErrZipCodeNotFound)
		rr := setup(cep, mockSvc)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "can not find zipcode\n", rr.Body.String())
		mockSvc.AssertExpectations(t)
	})

	t.Run("should return 500 for other errors", func(t *testing.T) {
		mockSvc := new(MockOrchestratorService)
		cep := "99999999"

		mockSvc.On("GetWeather", cep).Return(nil, errors.New("unexpected error"))
		rr := setup(cep, mockSvc)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "internal server error\n", rr.Body.String())
		mockSvc.AssertExpectations(t)
	})
}
