package http

import (
	"time"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http/handler"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http/middleware"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// RouterConfig holds configuration for the router
type RouterConfig struct {
	UserUsecase    domain.UserUsecase
	PostUsecase    domain.PostUsecase
	CommentUsecase domain.CommentUsecase
	LikeUsecase    domain.LikeUsecase
	Logger         *logrus.Logger
	JWTSecret      string
	AllowOrigins   []string
	RateLimit      float64
	RateBurst      int
}

// SetupRouter sets up the HTTP router with all handlers and middleware
func SetupRouter(config *RouterConfig) *gin.Engine {
	// Create router with default middleware
	router := gin.New()

	// Recovery middleware
	router.Use(gin.Recovery())

	// Logger middleware
	router.Use(middleware.Logger(config.Logger))

	// CORS middleware
	corsConfig := &middleware.CORSConfig{
		AllowOrigins: config.AllowOrigins,
	}
	router.Use(middleware.CORS(corsConfig))

	// Rate limiter middleware
	rateLimiter := middleware.NewRateLimiter(rate.Limit(config.RateLimit), config.RateBurst)

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(config.JWTSecret)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Public routes with IP-based rate limiting
		public := v1.Group("")
		public.Use(rateLimiter.RateLimit())
		{
			// User registration and login
			handler.NewUserHandler(public, config.UserUsecase, authMiddleware)
		}

		// Protected routes with user-based rate limiting
		protected := v1.Group("")
		protected.Use(authMiddleware.AuthRequired())
		protected.Use(rateLimiter.RateLimitByUser())
		{
			// Initialize handlers
			handler.NewPostHandler(protected, config.PostUsecase, authMiddleware)
			handler.NewCommentHandler(protected, config.CommentUsecase, authMiddleware)
			handler.NewLikeHandler(protected, config.LikeUsecase, authMiddleware)
		}
	}

	return router
}
