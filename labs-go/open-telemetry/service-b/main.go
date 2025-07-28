package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type ViaCepResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

type WeatherApiResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type CepInput struct {
	Cep string `json:"cep"`
}

var weatherApiKey string

func main() {
	weatherApiKey = os.Getenv("WEATHER_API_KEY")
	if weatherApiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable not set")
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	zipkinEndpoint := os.Getenv("OTEL_EXPORTER_ZIPKIN_ENDPOINT")

	tp, err := InitTracer(serviceName, zipkinEndpoint)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	mux := http.NewServeMux()
	handler := http.HandlerFunc(WeatherHandler)
	otelHandler := otelhttp.NewHandler(handler, "WeatherHandler")
	mux.Handle("/", otelHandler)

	log.Printf("%s is running on port 8081", serviceName)
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("could not listen on port 8081: %v", err)
	}
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("service-b-tracer")

	var input CepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	location, err := findCityByCep(ctx, tracer, input.Cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempC, err := findTemperatureByCity(ctx, tracer, location)
	if err != nil {
		http.Error(w, "Error fetching temperature", http.StatusInternalServerError)
		return
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	response := WeatherResponse{
		City:  location,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func findCityByCep(ctx context.Context, tracer trace.Tracer, cep string) (string, error) {
	ctx, span := tracer.Start(ctx, "findCityByCep-span")
	defer span.End()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var viaCepData ViaCepResponse
	if err := json.Unmarshal(body, &viaCepData); err != nil {
		return "", err
	}

	if viaCepData.Erro {
		return "", fmt.Errorf("cep not found")
	}

	return viaCepData.Localidade, nil
}

func findTemperatureByCity(ctx context.Context, tracer trace.Tracer, city string) (float64, error) {
	ctx, span := tracer.Start(ctx, "findTemperatureByCity-span")
	defer span.End()

	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", weatherApiKey, encodedCity)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var weatherData WeatherApiResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return 0, err
	}

	return weatherData.Current.TempC, nil
}

func InitTracer(serviceName, collectorURL string) (*sdktrace.TracerProvider, error) {
	exporter, err := zipkin.New(collectorURL)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
