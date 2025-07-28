package limiter

import (
	"context"
	"fmt"

	"rate-limiter/internal/config"
	"rate-limiter/internal/storage"
)

type RateLimiter struct {
	storage storage.LimiterStorage
	config  *config.Config
}

func NewRateLimiter(storage storage.LimiterStorage, config *config.Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, identifier, token string) (bool, error) {
	var key string
	var limit int64
	var blockDuration int64

	tokenLimit, isTokenConfigured := rl.config.TokenLimits[token]

	if token != "" && isTokenConfigured {
		key = fmt.Sprintf("token:%s", token)
		limit = tokenLimit
		blockDuration = rl.config.TokenBlockDurationSeconds
	} else {
		key = fmt.Sprintf("ip:%s", identifier)
		limit = rl.config.IPRequestsPerSecond
		blockDuration = rl.config.IPBlockDurationSeconds
	}

	blocked, err := rl.storage.IsBlocked(ctx, key)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	allowed, err := rl.storage.Allow(ctx, key, limit, 1) // 1 second window
	if err != nil {
		return false, err
	}

	if !allowed {
		if err := rl.storage.Block(ctx, key, blockDuration); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}