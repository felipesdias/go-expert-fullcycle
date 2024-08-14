package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type CepResult struct {
	Service      string
	Cep          string
	State        string
	City         string
	Neighborhood string
	Street       string
}

type CepProvider interface {
	getName() string
	findCep(string) (*CepResult, error)
}

func main() {
	cep := "36576202"

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	services := []CepProvider{
		NewBrasilApi(ctx),
		NewViaCep(ctx),
	}

	c := make(chan *CepResult)

	for _, service := range services {
		go func() {
			cep, err := service.findCep(cep)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Println("Error for", service.getName(), err)
			} else {
				c <- cep
				cancel()
			}
		}()
	}

	select {
	case cepResult := <-c:
		fmt.Printf("Service: %v\nAddress: %v, %v, %v - %v (%v)\n", cepResult.Service, cepResult.Street, cepResult.Neighborhood, cepResult.City, cepResult.State, cepResult.Cep)
	case <-ctx.Done():
		log.Fatal("Timeout")
	}
}
