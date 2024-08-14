package main

import "context"

type ViaCepDto struct {
	Cep          string `json:"cep"`
	Street       string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
}

type ViaCep struct {
	ctx  context.Context
	Name string
}

func NewViaCep(ctx context.Context) *ViaCep {
	return &ViaCep{
		Name: "ViaCep",
		ctx:  ctx,
	}
}

func (a ViaCep) getName() string {
	return a.Name
}

func (a ViaCep) findCep(cep string) (*CepResult, error) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	rDto, err := doGet[ViaCepDto](url, a.ctx)

	if err != nil {
		return nil, err
	}

	return &CepResult{
		Service:      a.Name,
		Cep:          rDto.Cep,
		City:         rDto.City,
		State:        rDto.State,
		Neighborhood: rDto.Neighborhood,
		Street:       rDto.Street,
	}, nil
}
