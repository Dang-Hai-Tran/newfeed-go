package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/repository/cache"
)

type postUsecase struct {
	postRepo    domain.PostRepository
	postCache   cache.PostCache
	userRepo    domain.UserRepository
	contextTimeout time.Duration
}

// NewPostUsecase creates a new post usecase
func NewPostUsecase(pr domain.PostRepository, pc cache.PostCache, ur domain.UserRepository, timeout time.Duration) domain.PostUsecase {
	return &postUsecase{
		postRepo:    pr,
		postCache:   pc,
		userRepo:    ur,
		contextTimeout: timeout,
	}
}

func (p *postUsecase) CreatePost(post *domain.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.contextTimeout)
	defer cancel()

	// Verify user exists
	user, err := p.userRepo.GetByID(post.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Set timestamps
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	// Create post in database
	if err := p.postRepo.Create(post); err != nil {
		return err
	}

	// Cache post
	if err := p.postCache.SetPost(ctx, post); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	// Invalidate user's posts cache and newsfeed cache for followers
	pattern := fmt.Sprintf("user:%d:posts:*", post.UserID)
	if err := p.postCache.DeletePattern(ctx, pattern); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (p *postUsecase) GetPost(id uint64) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.contextTimeout)
	defer cancel()

	// Try to get from cache first
	post, err := p.postCache.GetPost(ctx, id)
	if err == nil && post != nil {
		return post, nil
	}

	// If not in cache, get from database
	post, err = p.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// Cache post
	if err := p.postCache.SetPost(ctx, post); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return post, nil
}

func (p *postUsecase) GetUserPosts(userID uint64, page, limit int) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.contextTimeout)
	defer cancel()

	// Try to get from cache first
	posts, err := p.postCache.GetUserPosts(ctx, userID, page)
	if err == nil {
		return posts, nil
	}

	// If not in cache, get from database
	posts, err = p.postRepo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, err
	}

	// Cache posts
	if err := p.postCache.SetUserPosts(ctx, userID, page, posts); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return posts, nil
}

func (p *postUsecase) UpdatePost(post *domain.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.contextTimeout)
	defer cancel()

	// Verify post exists and belongs to user
	existingPost, err := p.postRepo.GetByID(post.ID)
	if err != nil {
		return err
	}
	if existingPost == nil {
		return errors.New("post not found")
	}
	if existingPost.UserID != post.UserID {
		return errors.New("unauthorized")
	}

	// Update timestamp
	post.UpdatedAt = time.Now()

	// Update in database
	if err := p.postRepo.Update(post); err != nil {
		return err
	}

	// Update in cache
	if err := p.postCache.SetPost(ctx, post); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	// Invalidate user's posts cache and newsfeed cache
	pattern := fmt.Sprintf("user:%d:posts:*", post.UserID)
	if err := p.postCache.DeletePattern(ctx, pattern); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (p *postUsecase) DeletePost(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.contextTimeout)
	defer cancel()

	// Delete from database
	if err := p.postRepo.Delete(id); err != nil {
		return err
	}

	// Delete from cache
	if err := p.postCache.DeletePost(ctx, id); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (p *postUsecase) GetNewsFeed(userID uint64, page, limit int) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.contextTimeout)
	defer cancel()

	// Try to get from cache first
	posts, err := p.postCache.GetNewsFeed(ctx, userID, page)
	if err == nil {
		return posts, nil
	}

	// If not in cache, get from database
	posts, err = p.postRepo.GetNewsFeed(userID, page, limit)
	if err != nil {
		return nil, err
	}

	// Cache newsfeed
	if err := p.postCache.SetNewsFeed(ctx, userID, page, posts); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return posts, nil
}
