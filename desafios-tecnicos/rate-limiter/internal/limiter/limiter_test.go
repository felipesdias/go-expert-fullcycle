package limiter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rate-limiter/internal/config"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) Block(ctx context.Context, key string, blockDuration int64) error {
	args := m.Called(ctx, key, blockDuration)
	return args.Error(0)
}

func (m *MockStorage) Allow(ctx context.Context, key string, limit int64, window int64) (bool, error) {
	args := m.Called(ctx, key, limit, window)
	return args.Bool(0), args.Error(1)
}

func TestRateLimiter(t *testing.T) {
	cfg := &config.Config{
		IPRequestsPerSecond:       5,
		IPBlockDurationSeconds:    300,
		TokenBlockDurationSeconds: 600,
		TokenLimits: map[string]int64{
			"test-token": 10,
		},
	}
	ctx := context.Background()

	t.Run("should allow request for IP under limit", func(t *testing.T) {
		storageMock := new(MockStorage)
		limiter := NewRateLimiter(storageMock, cfg)

		storageMock.On("IsBlocked", ctx, "ip:127.0.0.1").Return(false, nil)
		storageMock.On("Allow", ctx, "ip:127.0.0.1", int64(5), int64(1)).Return(true, nil)

		allowed, err := limiter.Allow(ctx, "127.0.0.1", "")
		assert.NoError(t, err)
		assert.True(t, allowed)
		storageMock.AssertExpectations(t)
	})

	t.Run("should block request for IP over limit", func(t *testing.T) {
		storageMock := new(MockStorage)
		limiter := NewRateLimiter(storageMock, cfg)

		storageMock.On("IsBlocked", ctx, "ip:127.0.0.1").Return(false, nil)
		storageMock.On("Allow", ctx, "ip:127.0.0.1", int64(5), int64(1)).Return(false, nil)
		storageMock.On("Block", ctx, "ip:127.0.0.1", int64(300)).Return(nil)

		allowed, err := limiter.Allow(ctx, "127.0.0.1", "")
		assert.NoError(t, err)
		assert.False(t, allowed)
		storageMock.AssertExpectations(t)
	})

	t.Run("should block request for already blocked IP", func(t *testing.T) {
		storageMock := new(MockStorage)
		limiter := NewRateLimiter(storageMock, cfg)

		storageMock.On("IsBlocked", ctx, "ip:127.0.0.1").Return(true, nil)

		allowed, err := limiter.Allow(ctx, "127.0.0.1", "")
		assert.NoError(t, err)
		assert.False(t, allowed)
		storageMock.AssertNotCalled(t, "Allow")
		storageMock.AssertNotCalled(t, "Block")
	})

	t.Run("should allow request for token under limit", func(t *testing.T) {
		storageMock := new(MockStorage)
		limiter := NewRateLimiter(storageMock, cfg)

		storageMock.On("IsBlocked", ctx, "token:test-token").Return(false, nil)
		storageMock.On("Allow", ctx, "token:test-token", int64(10), int64(1)).Return(true, nil)

		allowed, err := limiter.Allow(ctx, "127.0.0.1", "test-token")
		assert.NoError(t, err)
		assert.True(t, allowed)
		storageMock.AssertExpectations(t)
	})

	t.Run("should block request for token over limit", func(t *testing.T) {
		storageMock := new(MockStorage)
		limiter := NewRateLimiter(storageMock, cfg)

		storageMock.On("IsBlocked", ctx, "token:test-token").Return(false, nil)
		storageMock.On("Allow", ctx, "token:test-token", int64(10), int64(1)).Return(false, nil)
		storageMock.On("Block", ctx, "token:test-token", int64(600)).Return(nil)

		allowed, err := limiter.Allow(ctx, "127.0.0.1", "test-token")
		assert.NoError(t, err)
		assert.False(t, allowed)
		storageMock.AssertExpectations(t)
	})

	t.Run("should use IP limit for non-configured token", func(t *testing.T) {
		storageMock := new(MockStorage)
		limiter := NewRateLimiter(storageMock, cfg)

		storageMock.On("IsBlocked", ctx, "ip:127.0.0.1").Return(false, nil)
		storageMock.On("Allow", ctx, "ip:127.0.0.1", int64(5), int64(1)).Return(true, nil)

		allowed, err := limiter.Allow(ctx, "127.0.0.1", "unconfigured-token")
		assert.NoError(t, err)
		assert.True(t, allowed)
		storageMock.AssertExpectations(t)
	})
	
	t.Run("should handle storage error on IsBlocked", func(t *testing.T) {
		storageMock := new(MockStorage)
		limiter := NewRateLimiter(storageMock, cfg)
		expectedErr := errors.New("storage error")

		storageMock.On("IsBlocked", ctx, "ip:127.0.0.1").Return(false, expectedErr)

		allowed, err := limiter.Allow(ctx, "127.0.0.1", "")
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.False(t, allowed)
	})
}