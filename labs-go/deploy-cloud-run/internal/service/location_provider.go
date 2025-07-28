package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"weather-api/internal/entity"
)

type LocationProvider interface {
	GetLocation(cep string) (*entity.ViaCepOutput, error)
}

type ViaCepProvider struct {
	BaseURL string
}

func NewViaCepProvider() *ViaCepProvider {
	return &ViaCepProvider{
		BaseURL: "https://viacep.com.br/ws",
	}
}

func (p *ViaCepProvider) GetLocation(cep string) (*entity.ViaCepOutput, error) {
	url := fmt.Sprintf("%s/%s/json/", p.BaseURL, cep)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data entity.ViaCepOutput
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.IsErro() {
		return nil, nil // CEP válido, mas não encontrado
	}

	return &data, nil
}
