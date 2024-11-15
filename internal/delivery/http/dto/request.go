package dto

import "time"

type RegisterRequest struct {
	Username  string    `json:"username" binding:"required,min=3,max=50"`
	Password  string    `json:"password" binding:"required,min=6"`
	Email     string    `json:"email" binding:"required,email"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Birthday  time.Time `json:"birthday" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	Password  string    `json:"password"`
}

type CreatePostRequest struct {
	Content  string `json:"content" binding:"required"`
	ImageURL string `json:"image_url"`
}

type UpdatePostRequest struct {
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type PaginationQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}
