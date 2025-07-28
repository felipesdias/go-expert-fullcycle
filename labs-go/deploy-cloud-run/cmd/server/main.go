package main

import (
	"log"
	"net/http"
	"os"

	"weather-api/internal/handler"
	"weather-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable not set")
	}

	locationProvider := service.NewViaCepProvider()
	weatherProvider := service.NewWeatherApiProvider(apiKey)
	orchestratorService := service.NewOrchestratorService(locationProvider, weatherProvider)
	weatherHandler := handler.NewWeatherHandler(orchestratorService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/weather/{cep}", weatherHandler.GetWeather)

	log.Println("Server starting on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
