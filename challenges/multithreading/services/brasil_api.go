package services

import (
	"challenges_multithreading/models"
	"challenges_multithreading/utils"
	"context"
	"fmt"
)

type BrasilApiDto struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type BrasilApi struct {
	ctx context.Context
}

func NewBrasilApi(ctx context.Context) *BrasilApi {
	return &BrasilApi{ctx: ctx}
}

func (b *BrasilApi) GetName() string {
	return "BrasilApi"
}

func (b *BrasilApi) FindCep(cep string) (*models.CepResult, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	brasilApiDto, err := utils.DoGet[BrasilApiDto](url, b.ctx)
	if err != nil {
		return nil, err
	}

	return &models.CepResult{
		Service:      fmt.Sprintf("%s - %s", b.GetName(), brasilApiDto.Service),
		Cep:          brasilApiDto.Cep,
		City:         brasilApiDto.City,
		State:        brasilApiDto.State,
		Neighborhood: brasilApiDto.Neighborhood,
		Street:       brasilApiDto.Street,
	}, nil
}
