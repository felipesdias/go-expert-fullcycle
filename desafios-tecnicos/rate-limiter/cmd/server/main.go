package main

import (
	"log"
	"net/http"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	redisStorage := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	rateLimiter := limiter.NewRateLimiter(redisStorage, cfg)
	rlMiddleware := middleware.RateLimiterMiddleware(rateLimiter)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", rlMiddleware(mux)); err != nil {
		log.Fatalf("could not listen on port 8080: %v\n", err)
	}
}