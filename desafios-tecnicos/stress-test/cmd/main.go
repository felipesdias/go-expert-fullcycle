package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Report struct {
	TotalTime     time.Duration
	TotalRequests int
	StatusOK      int
	StatusDist    map[int]int
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado.")
	requests := flag.Int("requests", 0, "Número total de requests.")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas (workers).")
	flag.Parse()

	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	startTime := time.Now()

	jobs := make(chan struct{}, *requests)
	resultsChan := make(chan int, *requests)
	var wg sync.WaitGroup

	// Inicia o número definido de workers
	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go worker(*url, &wg, jobs, resultsChan)
	}

	// Envia todas as tarefas para o canal de jobs
	for j := 0; j < *requests; j++ {
		jobs <- struct{}{}
	}
	close(jobs)

	// Aguarda a finalização de todos os workers
	wg.Wait()
	close(resultsChan)

	totalTime := time.Since(startTime)

	report := processResults(*requests, resultsChan, totalTime)
	printReport(report)
}

// Worker processa jobs de um canal até que o canal seja fechado.
func worker(url string, wg *sync.WaitGroup, jobs <-chan struct{}, results chan<- int) {
	defer wg.Done()
	for range jobs {
		resp, err := http.Get(url)
		if err != nil {
			results <- 0
			continue
		}
		results <- resp.StatusCode
		resp.Body.Close()
	}
}

func processResults(totalRequests int, results <-chan int, totalTime time.Duration) Report {
	report := Report{
		TotalTime:     totalTime,
		TotalRequests: totalRequests,
		StatusDist:    make(map[int]int),
	}

	for statusCode := range results {
		if statusCode == http.StatusOK {
			report.StatusOK++
		} else {
			report.StatusDist[statusCode]++
		}
	}
	return report
}

func printReport(r Report) {
	fmt.Println("Relatório de Teste de Carga")
	fmt.Println("---------------------------")
	fmt.Printf("Tempo total gasto: %s\n", r.TotalTime)
	fmt.Printf("Quantidade total de requests: %d\n", r.TotalRequests)
	fmt.Printf("Requests com status 200 (OK): %d\n", r.StatusOK)

	if len(r.StatusDist) > 0 {
		fmt.Println("\nDistribuição de outros status:")
		for code, count := range r.StatusDist {
			if code == 0 {
				fmt.Printf("  - Erros de requisição: %d\n", count)
				continue
			}
			fmt.Printf("  - Status %d (%s): %d\n", code, http.StatusText(code), count)
		}
	}
	fmt.Println("---------------------------")
}
