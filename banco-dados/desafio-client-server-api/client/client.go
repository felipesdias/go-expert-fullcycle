package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Erro ao fazer requisição para o servidor: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Servidor retornou erro: %s, status: %s", string(body), resp.Status)
	}

	var cotacao CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		log.Fatalf("Erro ao decodificar resposta JSON: %v", err)
	}

	err = saveCotacaoToFile(cotacao.Bid)
	if err != nil {
		log.Fatalf("Erro ao salvar cotação em arquivo: %v", err)
	}

	log.Printf("Cotação salva com sucesso! Valor: %s", cotacao.Bid)
}
func saveCotacaoToFile(bid string) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", bid))
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo: %w", err)
	}

	return nil
}
