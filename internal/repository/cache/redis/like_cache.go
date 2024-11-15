package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/datran42/newfeed-go/internal/domain"
	"github.com/datran42/newfeed-go/internal/repository/cache"
	redisClient "github.com/datran42/newfeed-go/pkg/cache"
)

type likeCache struct {
	redis *redisClient.RedisClient
}

// NewLikeCache creates a new Redis like cache
func NewLikeCache(redis *redisClient.RedisClient) cache.LikeCache {
	return &likeCache{redis: redis}
}

func (c *likeCache) GetPostLikes(ctx context.Context, postID uint64, page int) ([]domain.Like, error) {
	key := fmt.Sprintf("post:%d:likes:page:%d", postID, page)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var likes []domain.Like
	if err := json.Unmarshal([]byte(data), &likes); err != nil {
		return nil, err
	}

	return likes, nil
}

func (c *likeCache) SetPostLikes(ctx context.Context, postID uint64, page int, likes []domain.Like) error {
	key := fmt.Sprintf("post:%d:likes:page:%d", postID, page)
	data, err := json.Marshal(likes)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, cache.DefaultCacheDuration)
}

func (c *likeCache) DeletePostLikes(ctx context.Context, postID uint64) error {
	pattern := fmt.Sprintf("post:%d:likes:*", postID)
	return c.redis.DeletePattern(ctx, pattern)
}

func (c *likeCache) GetLikeExists(ctx context.Context, postID, userID uint64) (bool, error) {
	key := fmt.Sprintf("post:%d:like:%d", postID, userID)
	exists, err := c.redis.Get(ctx, key)
	if err != nil {
		return false, err
	}

	return exists == "1", nil
}

func (c *likeCache) SetLikeExists(ctx context.Context, postID, userID uint64, exists bool) error {
	key := fmt.Sprintf("post:%d:like:%d", postID, userID)
	value := "0"
	if exists {
		value = "1"
	}

	// Cache like status for a longer duration since it changes less frequently
	return c.redis.Set(ctx, key, value, cache.LongCacheDuration)
}
