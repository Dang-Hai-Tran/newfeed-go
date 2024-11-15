package postgres

import (
	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	// Create followers table for many-to-many relationship
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS followers (
			follower_id BIGINT NOT NULL,
			following_id BIGINT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (follower_id, following_id)
		)
	`).Error
	if err != nil {
		return err
	}

	// Create indexes for better query performance
	err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON followers(follower_id);
		CREATE INDEX IF NOT EXISTS idx_followers_following_id ON followers(following_id);
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
		CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
		CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id);
		CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);
	`).Error
	if err != nil {
		return err
	}

	// Auto-migrate the schema
	return db.AutoMigrate(
		&domain.User{},
		&domain.Post{},
		&domain.Comment{},
		&domain.Like{},
	)
}
