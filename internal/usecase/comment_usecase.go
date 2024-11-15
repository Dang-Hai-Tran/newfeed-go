package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/datran42/newfeed-go/internal/domain"
	"github.com/datran42/newfeed-go/internal/repository/cache"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
	commentCache cache.CommentCache
	postRepo    domain.PostRepository
	userRepo    domain.UserRepository
	contextTimeout time.Duration
}

// NewCommentUsecase creates a new comment usecase
func NewCommentUsecase(
	cr domain.CommentRepository,
	cc cache.CommentCache,
	pr domain.PostRepository,
	ur domain.UserRepository,
	timeout time.Duration,
) domain.CommentUsecase {
	return &commentUsecase{
		commentRepo: cr,
		commentCache: cc,
		postRepo:    pr,
		userRepo:    ur,
		contextTimeout: timeout,
	}
}

func (c *commentUsecase) CreateComment(comment *domain.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.contextTimeout)
	defer cancel()

	// Verify user exists
	user, err := c.userRepo.GetByID(comment.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify post exists
	post, err := c.postRepo.GetByID(comment.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post not found")
	}

	// Set timestamps
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	// Create comment in database
	if err := c.commentRepo.Create(comment); err != nil {
		return err
	}

	// Invalidate post comments cache
	if err := c.commentCache.DeletePostComments(ctx, comment.PostID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (c *commentUsecase) GetComment(id uint64) (*domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.contextTimeout)
	defer cancel()

	comment, err := c.commentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("comment not found")
	}

	return comment, nil
}

func (c *commentUsecase) GetPostComments(postID uint64, page, limit int) ([]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.contextTimeout)
	defer cancel()

	// Try to get from cache first
	comments, err := c.commentCache.GetPostComments(ctx, postID, page)
	if err == nil {
		return comments, nil
	}

	// If not in cache, get from database
	comments, err = c.commentRepo.GetByPostID(postID, page, limit)
	if err != nil {
		return nil, err
	}

	// Cache comments
	if err := c.commentCache.SetPostComments(ctx, postID, page, comments); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return comments, nil
}

func (c *commentUsecase) UpdateComment(comment *domain.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.contextTimeout)
	defer cancel()

	// Verify comment exists and belongs to user
	existingComment, err := c.commentRepo.GetByID(comment.ID)
	if err != nil {
		return err
	}
	if existingComment == nil {
		return errors.New("comment not found")
	}
	if existingComment.UserID != comment.UserID {
		return errors.New("unauthorized")
	}

	// Update timestamp
	comment.UpdatedAt = time.Now()

	// Update in database
	if err := c.commentRepo.Update(comment); err != nil {
		return err
	}

	// Invalidate post comments cache
	if err := c.commentCache.DeletePostComments(ctx, comment.PostID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}

func (c *commentUsecase) DeleteComment(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.contextTimeout)
	defer cancel()

	// Get comment to know which post's cache to invalidate
	comment, err := c.commentRepo.GetByID(id)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Delete from database
	if err := c.commentRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate post comments cache
	if err := c.commentCache.DeletePostComments(ctx, comment.PostID); err != nil {
		// Log error but don't return it
		// TODO: Add proper logging
	}

	return nil
}
