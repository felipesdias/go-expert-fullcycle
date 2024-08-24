package services

import (
	"challenges_multithreading/models"
	"challenges_multithreading/utils"
	"context"
	"fmt"
)

type ViaCepDto struct {
	Cep          string `json:"cep"`
	Street       string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
}

type ViaCep struct {
	ctx context.Context
}

func NewViaCep(ctx context.Context) *ViaCep {
	return &ViaCep{ctx: ctx}
}

func (v *ViaCep) GetName() string {
	return "ViaCep"
}

func (v *ViaCep) FindCep(cep string) (*models.CepResult, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	viaCepDto, err := utils.DoGet[ViaCepDto](url, v.ctx)
	if err != nil {
		return nil, err
	}

	return &models.CepResult{
		Service:      v.GetName(),
		Cep:          viaCepDto.Cep,
		City:         viaCepDto.City,
		State:        viaCepDto.State,
		Neighborhood: viaCepDto.Neighborhood,
		Street:       viaCepDto.Street,
	}, nil
}
