package config

import (
	"log"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisAddr                 string `envconfig:"REDIS_ADDR"`
	RedisPassword             string `envconfig:"REDIS_PASSWORD"`
	RedisDB                   int    `envconfig:"REDIS_DB"`
	IPRequestsPerSecond       int64  `envconfig:"IP_REQUESTS_PER_SECOND"`
	IPBlockDurationSeconds    int64  `envconfig:"IP_BLOCK_DURATION_SECONDS"`
	TokenLimitsRaw            string `envconfig:"TOKEN_LIMITS"`
	TokenBlockDurationSeconds int64  `envconfig:"TOKEN_BLOCK_DURATION_SECONDS"`
	TokenLimits               map[string]int64
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	cfg.TokenLimits = make(map[string]int64)
	if cfg.TokenLimitsRaw != "" {
		pairs := strings.Split(cfg.TokenLimitsRaw, ",")
		for _, pair := range pairs {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				tokenKey := strings.TrimSpace(parts[0])
				limitStr := strings.TrimSpace(parts[1])

				limit, err := strconv.ParseInt(limitStr, 10, 64)
				if err != nil {
					log.Printf("Invalid limit value for token %s: %v", tokenKey, err)
					continue
				}
				log.Printf("Setting token limit for %s: %d", tokenKey, limit)
				cfg.TokenLimits[tokenKey] = limit
			}
		}
	}

	return &cfg, nil
}