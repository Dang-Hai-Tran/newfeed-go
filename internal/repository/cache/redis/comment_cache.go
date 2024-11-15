package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/datran42/newfeed-go/internal/domain"
	"github.com/datran42/newfeed-go/internal/repository/cache"
	redisClient "github.com/datran42/newfeed-go/pkg/cache"
)

type commentCache struct {
	redis *redisClient.RedisClient
}

// NewCommentCache creates a new Redis comment cache
func NewCommentCache(redis *redisClient.RedisClient) cache.CommentCache {
	return &commentCache{redis: redis}
}

func (c *commentCache) GetPostComments(ctx context.Context, postID uint64, page int) ([]domain.Comment, error) {
	key := fmt.Sprintf("post:%d:comments:page:%d", postID, page)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var comments []domain.Comment
	if err := json.Unmarshal([]byte(data), &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *commentCache) SetPostComments(ctx context.Context, postID uint64, page int, comments []domain.Comment) error {
	key := fmt.Sprintf("post:%d:comments:page:%d", postID, page)
	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, data, cache.DefaultCacheDuration)
}

func (c *commentCache) DeletePostComments(ctx context.Context, postID uint64) error {
	pattern := fmt.Sprintf("post:%d:comments:*", postID)
	return c.redis.DeletePattern(ctx, pattern)
}
