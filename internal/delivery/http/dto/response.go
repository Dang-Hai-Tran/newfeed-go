package dto

import (
	"time"

	"github.com/datran42/newfeed-go/internal/domain"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type UserResponse struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
}

type PostResponse struct {
	ID        uint64           `json:"id"`
	UserID    uint64          `json:"user_id"`
	Content   string          `json:"content"`
	ImageURL  string          `json:"image_url,omitempty"`
	Likes     []LikeResponse  `json:"likes,omitempty"`
	Comments  []CommentResponse `json:"comments,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CommentResponse struct {
	ID        uint64    `json:"id"`
	PostID    uint64    `json:"post_id"`
	UserID    uint64    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LikeResponse struct {
	ID        uint64    `json:"id"`
	PostID    uint64    `json:"post_id"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Convert domain models to response DTOs
func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Birthday:  user.Birthday,
		CreatedAt: user.CreatedAt,
	}
}

func ToPostResponse(post *domain.Post) *PostResponse {
	likes := make([]LikeResponse, len(post.Likes))
	for i, like := range post.Likes {
		likes[i] = *ToLikeResponse(&like)
	}

	comments := make([]CommentResponse, len(post.Comments))
	for i, comment := range post.Comments {
		comments[i] = *ToCommentResponse(&comment)
	}

	return &PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Content:   post.Content,
		ImageURL:  post.ImageURL,
		Likes:     likes,
		Comments:  comments,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func ToCommentResponse(comment *domain.Comment) *CommentResponse {
	return &CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func ToLikeResponse(like *domain.Like) *LikeResponse {
	return &LikeResponse{
		ID:        like.ID,
		PostID:    like.PostID,
		UserID:    like.UserID,
		CreatedAt: like.CreatedAt,
	}
}
