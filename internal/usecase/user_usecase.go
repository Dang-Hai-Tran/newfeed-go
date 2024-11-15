package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/repository/cache"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo    domain.UserRepository
	userCache   cache.UserCache
	contextTimeout time.Duration
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(ur domain.UserRepository, uc cache.UserCache, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:    ur,
		userCache:   uc,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Register(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Check if username exists
	existingUser, err := u.userRepo.GetByUsername(user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Check if email exists
	existingUser, err = u.userRepo.GetByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create user
	if err := u.userRepo.Create(user); err != nil {
		return err
	}

	// Cache user
	return u.userCache.SetUser(ctx, user)
}

func (u *userUsecase) Login(username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Get user from database
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid username or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// TODO: Generate JWT token
	// For now, return a placeholder
	return "jwt_token", nil
}

func (u *userUsecase) GetProfile(id uint64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Try to get from cache first
	user, err := u.userCache.GetUser(ctx, id)
	if err == nil && user != nil {
		return user, nil
	}

	// If not in cache, get from database
	user, err = u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Cache the user
	if err := u.userCache.SetUser(ctx, user); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return user, nil
}

func (u *userUsecase) UpdateProfile(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Update timestamp
	user.UpdatedAt = time.Now()

	// If password is being updated, hash it
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	// Update in database
	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	// Update in cache
	return u.userCache.SetUser(ctx, user)
}

func (u *userUsecase) DeleteProfile(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Delete from database
	if err := u.userRepo.Delete(id); err != nil {
		return err
	}

	// Delete from cache
	return u.userCache.DeleteUser(ctx, id)
}

func (u *userUsecase) Follow(followerID, followingID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Add follower in database
	if err := u.userRepo.Follow(followerID, followingID); err != nil {
		return err
	}

	// Invalidate followers and following cache
	if err := u.userCache.DeleteUser(ctx, followerID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}
	if err := u.userCache.DeleteUser(ctx, followingID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (u *userUsecase) Unfollow(followerID, followingID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Remove follower in database
	if err := u.userRepo.Unfollow(followerID, followingID); err != nil {
		return err
	}

	// Invalidate followers and following cache
	if err := u.userCache.DeleteUser(ctx, followerID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}
	if err := u.userCache.DeleteUser(ctx, followingID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (u *userUsecase) GetFollowers(userID uint64) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Try to get from cache first
	followers, err := u.userCache.GetFollowers(ctx, userID)
	if err == nil {
		return followers, nil
	}

	// If not in cache, get from database
	followers, err = u.userRepo.GetFollowers(userID)
	if err != nil {
		return nil, err
	}

	// Cache the followers
	if err := u.userCache.SetFollowers(ctx, userID, followers); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return followers, nil
}

func (u *userUsecase) GetFollowing(userID uint64) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()

	// Try to get from cache first
	following, err := u.userCache.GetFollowing(ctx, userID)
	if err == nil {
		return following, nil
	}

	// If not in cache, get from database
	following, err = u.userRepo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}

	// Cache the following
	if err := u.userCache.SetFollowing(ctx, userID, following); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return following, nil
}
