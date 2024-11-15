package domain

import (
	"time"
)

type Like struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	PostID    uint64    `json:"post_id" gorm:"not null"`
	UserID    uint64    `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeRepository interface {
	Create(like *Like) error
	Delete(postID, userID uint64) error
	GetByPostID(postID uint64, page, limit int) ([]Like, error)
	Exists(postID, userID uint64) (bool, error)
}

type LikeUsecase interface {
	LikePost(postID, userID uint64) error
	UnlikePost(postID, userID uint64) error
	GetPostLikes(postID uint64, page, limit int) ([]Like, error)
	HasUserLiked(postID, userID uint64) (bool, error)
}
