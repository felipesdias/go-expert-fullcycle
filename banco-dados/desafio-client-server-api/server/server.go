package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CotacaoAPIResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {

	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Falha na configuração do banco de dados: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", cotacaoHandler(db))

	log.Println("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacao (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func cotacaoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cotacao, err := getCotacaoFromAPI(r.Context())
		if err != nil {
			log.Printf("Erro ao buscar cotação da API: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = saveCotacaoInDB(db, cotacao.USDBRL.Bid)
		if err != nil {

			log.Printf("Erro ao salvar cotação no banco de dados: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(CotacaoResponse{Bid: cotacao.USDBRL.Bid})
		if err != nil {
			log.Printf("Erro ao encodar resposta para o cliente: %v", err)
			http.Error(w, "Erro interno", http.StatusInternalServerError)
		}
	}
}

func getCotacaoFromAPI(requestContext context.Context) (*CotacaoAPIResponse, error) {
	ctx, cancel := context.WithTimeout(requestContext, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição para a API de cotação: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API de cotação retornou status inesperado: %s", resp.Status)
	}

	var cotacao CotacaoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON da API de cotação: %w", err)
	}

	return &cotacao, nil
}

func saveCotacaoInDB(db *sql.DB, bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.Prepare("INSERT INTO cotacao(bid) VALUES(?)")
	if err != nil {
		return fmt.Errorf("erro ao preparar statement SQL: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid)
	if err != nil {
		return fmt.Errorf("erro ao executar insert no banco: %w", err)
	}

	log.Println("Cotação salva no banco com sucesso.")
	return nil
}
