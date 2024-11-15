package domain

import (
	"time"
)

type Post struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id" gorm:"not null"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	Likes     []Like    `json:"likes,omitempty" gorm:"foreignKey:PostID"`
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostRepository interface {
	Create(post *Post) error
	GetByID(id uint64) (*Post, error)
	GetByUserID(userID uint64, page, limit int) ([]Post, error)
	Update(post *Post) error
	Delete(id uint64) error
	GetNewsFeed(userID uint64, page, limit int) ([]Post, error)
}

type PostUsecase interface {
	CreatePost(post *Post) error
	GetPost(id uint64) (*Post, error)
	GetUserPosts(userID uint64, page, limit int) ([]Post, error)
	UpdatePost(post *Post) error
	DeletePost(id uint64) error
	GetNewsFeed(userID uint64, page, limit int) ([]Post, error)
}
