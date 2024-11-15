package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Config holds server configuration
type Config struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxHeaderBytes  int
	GracefulTimeout time.Duration
}

// DefaultConfig returns default server configuration
func DefaultConfig() *Config {
	return &Config{
		Host:            "0.0.0.0",
		Port:            8080,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		MaxHeaderBytes:  1 << 20, // 1 MB
		GracefulTimeout: 5 * time.Second,
	}
}

// Server represents HTTP server
type Server struct {
	server  *http.Server
	logger  *logrus.Logger
	config  *Config
	handler *gin.Engine
}

// NewServer creates a new HTTP server
func NewServer(handler *gin.Engine, logger *logrus.Logger, config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &Server{
		server:  server,
		logger:  logger,
		config:  config,
		handler: handler,
	}
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
	// Channel to receive error from server
	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		s.logger.Infof("Server is starting on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start server: %v", err)
		}
	}()

	// Channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		s.logger.Infof("Received signal: %v", sig)
		return s.Shutdown()
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.logger.Info("Server is shutting down...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), s.config.GracefulTimeout)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %v", err)
	}

	s.logger.Info("Server stopped gracefully")
	return nil
}

// Health returns a simple health check handler
func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	}
}
