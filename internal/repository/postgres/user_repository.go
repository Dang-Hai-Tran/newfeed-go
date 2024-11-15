package postgres

import (
	"errors"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint64) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *userRepository) GetFollowers(userID uint64) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Raw(`
		SELECT u.* FROM users u
		INNER JOIN followers f ON f.follower_id = u.id
		WHERE f.following_id = ?
	`, userID).Scan(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetFollowing(userID uint64) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Raw(`
		SELECT u.* FROM users u
		INNER JOIN followers f ON f.following_id = u.id
		WHERE f.follower_id = ?
	`, userID).Scan(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Follow(followerID, followingID uint64) error {
	return r.db.Exec(`
		INSERT INTO followers (follower_id, following_id)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING
	`, followerID, followingID).Error
}

func (r *userRepository) Unfollow(followerID, followingID uint64) error {
	return r.db.Exec(`
		DELETE FROM followers
		WHERE follower_id = ? AND following_id = ?
	`, followerID, followingID).Error
}
