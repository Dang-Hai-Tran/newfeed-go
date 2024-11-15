package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/repository/cache"
)

type likeUsecase struct {
	likeRepo    domain.LikeRepository
	likeCache   cache.LikeCache
	postRepo    domain.PostRepository
	userRepo    domain.UserRepository
	contextTimeout time.Duration
}

// NewLikeUsecase creates a new like usecase
func NewLikeUsecase(
	lr domain.LikeRepository,
	lc cache.LikeCache,
	pr domain.PostRepository,
	ur domain.UserRepository,
	timeout time.Duration,
) domain.LikeUsecase {
	return &likeUsecase{
		likeRepo:    lr,
		likeCache:   lc,
		postRepo:    pr,
		userRepo:    ur,
		contextTimeout: timeout,
	}
}

func (l *likeUsecase) LikePost(postID, userID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), l.contextTimeout)
	defer cancel()

	// Verify user exists
	user, err := l.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify post exists
	post, err := l.postRepo.GetByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post not found")
	}

	// Check if already liked
	exists, err := l.likeRepo.Exists(postID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("post already liked")
	}

	// Create like
	like := &domain.Like{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := l.likeRepo.Create(like); err != nil {
		return err
	}

	// Update cache
	if err := l.likeCache.SetLikeExists(ctx, postID, userID, true); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	// Invalidate post likes cache
	if err := l.likeCache.DeletePostLikes(ctx, postID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (l *likeUsecase) UnlikePost(postID, userID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), l.contextTimeout)
	defer cancel()

	// Check if like exists
	exists, err := l.likeRepo.Exists(postID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("post not liked")
	}

	// Delete like
	if err := l.likeRepo.Delete(postID, userID); err != nil {
		return err
	}

	// Update cache
	if err := l.likeCache.SetLikeExists(ctx, postID, userID, false); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	// Invalidate post likes cache
	if err := l.likeCache.DeletePostLikes(ctx, postID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (l *likeUsecase) GetPostLikes(postID uint64, page, limit int) ([]domain.Like, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.contextTimeout)
	defer cancel()

	// Try to get from cache first
	likes, err := l.likeCache.GetPostLikes(ctx, postID, page)
	if err == nil {
		return likes, nil
	}

	// If not in cache, get from database
	likes, err = l.likeRepo.GetByPostID(postID, page, limit)
	if err != nil {
		return nil, err
	}

	// Cache likes
	if err := l.likeCache.SetPostLikes(ctx, postID, page, likes); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return likes, nil
}

func (l *likeUsecase) HasUserLiked(postID, userID uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.contextTimeout)
	defer cancel()

	// Try to get from cache first
	exists, err := l.likeCache.GetLikeExists(ctx, postID, userID)
	if err == nil {
		return exists, nil
	}

	// If not in cache, get from database
	exists, err = l.likeRepo.Exists(postID, userID)
	if err != nil {
		return false, err
	}

	// Cache the result
	if err := l.likeCache.SetLikeExists(ctx, postID, userID, exists); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return exists, nil
}
