package main

import "context"

type brasilApiDto struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type BrasilApi struct {
	ctx  context.Context
	Name string
}

func NewBrasilApi(ctx context.Context) *BrasilApi {
	return &BrasilApi{
		Name: "BrasilApi",
		ctx:  ctx,
	}
}

func (a BrasilApi) getName() string {
	return a.Name
}

func (a BrasilApi) findCep(cep string) (*CepResult, error) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	rDto, err := doGet[brasilApiDto](url, a.ctx)

	if err != nil {
		return nil, err
	}

	return &CepResult{
		Service:      a.Name + " - " + rDto.Service,
		Cep:          rDto.Cep,
		City:         rDto.City,
		State:        rDto.State,
		Neighborhood: rDto.Neighborhood,
		Street:       rDto.Street,
	}, nil
}
