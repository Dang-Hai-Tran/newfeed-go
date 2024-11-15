package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/datran42/newfeed-go/config"
	"github.com/go-redis/redis/v8"
)

// RedisClient wraps redis.Client with additional functionality
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client instance
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxRetries:   3,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: client}, nil
}

// Set stores a key-value pair with expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching the pattern
func (r *RedisClient) DeletePattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

// SetHash stores a hash map with expiration
func (r *RedisClient) SetHash(ctx context.Context, key string, values map[string]interface{}, expiration time.Duration) error {
	pipe := r.client.Pipeline()
	pipe.HMSet(ctx, key, values)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)
	return err
}

// GetHash retrieves all fields from a hash
func (r *RedisClient) GetHash(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}
