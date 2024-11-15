package domain

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uint64) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint64) error
	GetFollowers(userID uint64) ([]User, error)
	GetFollowing(userID uint64) ([]User, error)
	Follow(followerID, followingID uint64) error
	Unfollow(followerID, followingID uint64) error
}

type UserUsecase interface {
	Register(user *User) error
	Login(username, password string) (string, error)
	GetProfile(id uint64) (*User, error)
	UpdateProfile(user *User) error
	DeleteProfile(id uint64) error
	Follow(followerID, followingID uint64) error
	Unfollow(followerID, followingID uint64) error
	GetFollowers(userID uint64) ([]User, error)
	GetFollowing(userID uint64) ([]User, error)
}
