package main

import (
	"log"
	"time"

	"github.com/Dang-Hai-Tran/newfeed-go/config"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http/server"
	"github.com/Dang-Hai-Tran/newfeed-go/pkg/cache"
	"github.com/Dang-Hai-Tran/newfeed-go/pkg/database"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis cache
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	userCache := redis.NewUserCache(redisClient)
	postRepo := postgres.NewPostRepository(db)
	postCache := redis.NewPostCache(redisClient)
	commentRepo := postgres.NewCommentRepository(db)
	commentCache := redis.NewCommentCache(redisClient)
	likeRepo := postgres.NewLikeRepository(db)
	likeCache := redis.NewLikeCache(redisClient)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, userCache, cfg.Auth.TokenExpiry)
	postUsecase := usecase.NewPostUsecase(postRepo, postCache, userRepo, cfg.Context.Timeout)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, commentCache, postRepo, userRepo, cfg.Context.Timeout)
	likeUsecase := usecase.NewLikeUsecase(likeRepo, likeCache, postRepo, userRepo, cfg.Context.Timeout)

	// Setup router
	routerConfig := &http.RouterConfig{
		UserUsecase:    userUsecase,
		PostUsecase:    postUsecase,
		CommentUsecase: commentUsecase,
		LikeUsecase:    likeUsecase,
		Logger:         logger,
		JWTSecret:      cfg.Auth.JWTSecret,
		AllowOrigins:   cfg.CORS.AllowOrigins,
		RateLimit:      cfg.RateLimit.Rate,
		RateBurst:      cfg.RateLimit.Burst,
	}
	router := http.SetupRouter(routerConfig)

	// Add health check endpoint
	router.GET("/health", server.Health())

	// Setup and start server
	serverConfig := &server.Config{
		Host:            cfg.Server.Host,
		Port:            cfg.Server.Port,
		ReadTimeout:     time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes:  cfg.Server.MaxHeaderBytes,
		GracefulTimeout: time.Duration(cfg.Server.GracefulTimeout) * time.Second,
	}

	srv := server.NewServer(router, logger, serverConfig)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
