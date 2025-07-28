package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"

	"rate-limiter/internal/limiter"
)

func RateLimiterMiddleware(limiter *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("API_KEY")

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
				if strings.Contains(ip, "[::1]") {
					ip = "127.0.0.1"
				}
			}

			allowed, err := limiter.Allow(r.Context(), ip, token)
			if err != nil {
				log.Printf("Error checking rate limit: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}