package postgres

import (
	"errors"

	"github.com/datran42/newfeed-go/internal/domain"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new instance of CommentRepository
func NewCommentRepository(db *gorm.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByID(id uint64) (*domain.Comment, error) {
	var comment domain.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) GetByPostID(postID uint64, page, limit int) ([]domain.Comment, error) {
	var comments []domain.Comment
	offset := (page - 1) * limit

	err := r.db.Where("post_id = ?", postID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) Update(comment *domain.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Comment{}, id).Error
}
