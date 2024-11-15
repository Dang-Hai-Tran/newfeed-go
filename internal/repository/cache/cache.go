package cache

import (
	"context"
	"time"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
)

// Default cache durations
const (
	DefaultCacheDuration = 15 * time.Minute
	LongCacheDuration   = 1 * time.Hour
	ShortCacheDuration  = 5 * time.Minute
)

type UserCache interface {
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	SetUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint64) error
	GetFollowers(ctx context.Context, userID uint64) ([]domain.User, error)
	SetFollowers(ctx context.Context, userID uint64, followers []domain.User) error
	GetFollowing(ctx context.Context, userID uint64) ([]domain.User, error)
	SetFollowing(ctx context.Context, userID uint64, following []domain.User) error
}

type PostCache interface {
	GetPost(ctx context.Context, id uint64) (*domain.Post, error)
	SetPost(ctx context.Context, post *domain.Post) error
	DeletePost(ctx context.Context, id uint64) error
	GetUserPosts(ctx context.Context, userID uint64, page int) ([]domain.Post, error)
	SetUserPosts(ctx context.Context, userID uint64, page int, posts []domain.Post) error
	GetNewsFeed(ctx context.Context, userID uint64, page int) ([]domain.Post, error)
	SetNewsFeed(ctx context.Context, userID uint64, page int, posts []domain.Post) error
}

type CommentCache interface {
	GetPostComments(ctx context.Context, postID uint64, page int) ([]domain.Comment, error)
	SetPostComments(ctx context.Context, postID uint64, page int, comments []domain.Comment) error
	DeletePostComments(ctx context.Context, postID uint64) error
}

type LikeCache interface {
	GetPostLikes(ctx context.Context, postID uint64, page int) ([]domain.Like, error)
	SetPostLikes(ctx context.Context, postID uint64, page int, likes []domain.Like) error
	DeletePostLikes(ctx context.Context, postID uint64) error
	GetLikeExists(ctx context.Context, postID, userID uint64) (bool, error)
	SetLikeExists(ctx context.Context, postID, userID uint64, exists bool) error
}
