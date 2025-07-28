package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStorage{client: rdb}
}

func (s *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	val, err := s.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (s *RedisStorage) Block(ctx context.Context, key string, blockDuration int64) error {
	blockKey := fmt.Sprintf("block:%s", key)
	return s.client.Set(ctx, blockKey, "blocked", time.Duration(blockDuration)*time.Second).Err()
}

func (s *RedisStorage) Allow(ctx context.Context, key string, limit int64, window int64) (bool, error) {
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		s.client.Expire(ctx, key, time.Duration(window)*time.Second)
	}

	return count <= limit, nil
}