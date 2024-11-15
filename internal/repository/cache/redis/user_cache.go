package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/datran42/newfeed-go/internal/domain"
	"github.com/datran42/newfeed-go/internal/repository/cache"
	redisClient "github.com/datran42/newfeed-go/pkg/cache"
)

type userCache struct {
	redis *redisClient.RedisClient
}

// NewUserCache creates a new Redis user cache
func NewUserCache(redis *redisClient.RedisClient) cache.UserCache {
	return &userCache{redis: redis}
}

func (c *userCache) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	key := fmt.Sprintf("user:%d", id)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *userCache) SetUser(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, cache.DefaultCacheDuration)
}

func (c *userCache) DeleteUser(ctx context.Context, id uint64) error {
	key := fmt.Sprintf("user:%d", id)
	return c.redis.Delete(ctx, key)
}

func (c *userCache) GetFollowers(ctx context.Context, userID uint64) ([]domain.User, error) {
	key := fmt.Sprintf("user:%d:followers", userID)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var followers []domain.User
	if err := json.Unmarshal([]byte(data), &followers); err != nil {
		return nil, err
	}

	return followers, nil
}

func (c *userCache) SetFollowers(ctx context.Context, userID uint64, followers []domain.User) error {
	key := fmt.Sprintf("user:%d:followers", userID)
	data, err := json.Marshal(followers)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, cache.DefaultCacheDuration)
}

func (c *userCache) GetFollowing(ctx context.Context, userID uint64) ([]domain.User, error) {
	key := fmt.Sprintf("user:%d:following", userID)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var following []domain.User
	if err := json.Unmarshal([]byte(data), &following); err != nil {
		return nil, err
	}

	return following, nil
}

func (c *userCache) SetFollowing(ctx context.Context, userID uint64, following []domain.User) error {
	key := fmt.Sprintf("user:%d:following", userID)
	data, err := json.Marshal(following)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, cache.DefaultCacheDuration)
}
