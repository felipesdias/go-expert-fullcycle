package storage

import "context"

type LimiterStorage interface {
	Allow(ctx context.Context, key string, limit, window int64) (bool, error)
	Block(ctx context.Context, key string, blockDuration int64) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}