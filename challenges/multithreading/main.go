package main

import (
	"challenges_multithreading/models"
	"challenges_multithreading/services"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	cep := "36576202"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	servicesList := []services.CepProvider{
		services.NewBrasilApi(ctx),
		services.NewViaCep(ctx),
	}

	results := make(chan *models.CepResult)

	for _, service := range servicesList {
		go func(s services.CepProvider) {
			result, err := s.FindCep(cep)
			if err != nil && ctx.Err() == nil {
				log.Println("Error for", s.GetName(), err)
			} else {
				results <- result
				cancel()
			}
		}(service)
	}

	select {
	case result := <-results:
		fmt.Printf("Service: %v\nAddress: %v, %v, %v - %v (%v)\n", result.Service, result.Street, result.Neighborhood, result.City, result.State, result.Cep)
	case <-ctx.Done():
		log.Fatal("Timeout")
	}
}
