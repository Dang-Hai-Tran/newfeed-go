package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/datran42/newfeed-go/internal/domain"
	"github.com/datran42/newfeed-go/internal/repository/cache"
	redisClient "github.com/datran42/newfeed-go/pkg/cache"
)

type postCache struct {
	redis *redisClient.RedisClient
}

// NewPostCache creates a new Redis post cache
func NewPostCache(redis *redisClient.RedisClient) cache.PostCache {
	return &postCache{redis: redis}
}

func (c *postCache) GetPost(ctx context.Context, id uint64) (*domain.Post, error) {
	key := fmt.Sprintf("post:%d", id)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var post domain.Post
	if err := json.Unmarshal([]byte(data), &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (c *postCache) SetPost(ctx context.Context, post *domain.Post) error {
	key := fmt.Sprintf("post:%d", post.ID)
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, cache.DefaultCacheDuration)
}

func (c *postCache) DeletePost(ctx context.Context, id uint64) error {
	key := fmt.Sprintf("post:%d", id)
	pattern := fmt.Sprintf("post:%d:*", id)
	
	// Delete post and all related cached data
	if err := c.redis.Delete(ctx, key); err != nil {
		return err
	}
	return c.redis.DeletePattern(ctx, pattern)
}

func (c *postCache) GetUserPosts(ctx context.Context, userID uint64, page int) ([]domain.Post, error) {
	key := fmt.Sprintf("user:%d:posts:page:%d", userID, page)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var posts []domain.Post
	if err := json.Unmarshal([]byte(data), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *postCache) SetUserPosts(ctx context.Context, userID uint64, page int, posts []domain.Post) error {
	key := fmt.Sprintf("user:%d:posts:page:%d", userID, page)
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, cache.DefaultCacheDuration)
}

func (c *postCache) GetNewsFeed(ctx context.Context, userID uint64, page int) ([]domain.Post, error) {
	key := fmt.Sprintf("user:%d:newsfeed:page:%d", userID, page)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var posts []domain.Post
	if err := json.Unmarshal([]byte(data), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *postCache) SetNewsFeed(ctx context.Context, userID uint64, page int, posts []domain.Post) error {
	key := fmt.Sprintf("user:%d:newsfeed:page:%d", userID, page)
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	// Cache newsfeed for a shorter duration since it's more dynamic
	return c.redis.Set(ctx, key, data, cache.ShortCacheDuration)
}
