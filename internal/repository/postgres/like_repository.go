package postgres

import (
	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

// NewLikeRepository creates a new instance of LikeRepository
func NewLikeRepository(db *gorm.DB) domain.LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *domain.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(postID, userID uint64) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&domain.Like{}).Error
}

func (r *likeRepository) GetByPostID(postID uint64, page, limit int) ([]domain.Like, error) {
	var likes []domain.Like
	offset := (page - 1) * limit

	err := r.db.Where("post_id = ?", postID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error

	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *likeRepository) Exists(postID, userID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Like{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
