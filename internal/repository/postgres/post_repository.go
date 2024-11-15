package postgres

import (
	"errors"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new instance of PostRepository
func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id uint64) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.Preload("Likes").Preload("Comments").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetByUserID(userID uint64, page, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	offset := (page - 1) * limit

	err := r.db.Where("user_id = ?", userID).
		Preload("Likes").
		Preload("Comments").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint64) error {
	// Start a transaction to delete post and related data
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete likes
		if err := tx.Where("post_id = ?", id).Delete(&domain.Like{}).Error; err != nil {
			return err
		}

		// Delete comments
		if err := tx.Where("post_id = ?", id).Delete(&domain.Comment{}).Error; err != nil {
			return err
		}

		// Delete post
		if err := tx.Delete(&domain.Post{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *postRepository) GetNewsFeed(userID uint64, page, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	offset := (page - 1) * limit

	err := r.db.Raw(`
		SELECT DISTINCT p.* FROM posts p
		INNER JOIN followers f ON f.following_id = p.user_id
		WHERE f.follower_id = ?
		UNION
		SELECT p.* FROM posts p
		WHERE p.user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, userID, limit, offset).
		Preload("Likes").
		Preload("Comments").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}
