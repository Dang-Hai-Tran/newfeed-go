package domain

import (
	"time"
)

type Comment struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	PostID    uint64    `json:"post_id" gorm:"not null"`
	UserID    uint64    `json:"user_id" gorm:"not null"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentRepository interface {
	Create(comment *Comment) error
	GetByID(id uint64) (*Comment, error)
	GetByPostID(postID uint64, page, limit int) ([]Comment, error)
	Update(comment *Comment) error
	Delete(id uint64) error
}

type CommentUsecase interface {
	CreateComment(comment *Comment) error
	GetComment(id uint64) (*Comment, error)
	GetPostComments(postID uint64, page, limit int) ([]Comment, error)
	UpdateComment(comment *Comment) error
	DeleteComment(id uint64) error
}
