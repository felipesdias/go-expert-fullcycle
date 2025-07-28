package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"weather-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	Orchestrator service.Orchestrator
}

func NewWeatherHandler(orch service.Orchestrator) *WeatherHandler {
	return &WeatherHandler{
		Orchestrator: orch,
	}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	weather, err := h.Orchestrator.GetWeather(cep)

	if err != nil {
		var status int
		var message string

		switch {
		case errors.Is(err, service.ErrInvalidZipCode):
			status = http.StatusUnprocessableEntity
			message = service.ErrInvalidZipCode.Error()
		case errors.Is(err, service.ErrZipCodeNotFound):
			status = http.StatusNotFound
			message = service.ErrZipCodeNotFound.Error()
		default:
			status = http.StatusInternalServerError
			message = "internal server error"
		}

		http.Error(w, message, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}
