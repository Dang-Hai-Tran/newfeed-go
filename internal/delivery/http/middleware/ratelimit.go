package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       *sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter with specified rate and burst
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		mu:       &sync.RWMutex{},
		rate:     r,
		burst:    b,
	}
}

// getLimiter returns a rate limiter for the given key
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// cleanup removes old limiters periodically
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Hour)
		rl.mu.Lock()
		rl.limiters = make(map[string]*rate.Limiter)
		rl.mu.Unlock()
	}
}

// RateLimit middleware function
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	go rl.cleanup() // Start cleanup goroutine

	return func(c *gin.Context) {
		// Use IP address as key for rate limiting
		key := c.ClientIP()
		limiter := rl.getLimiter(key)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser middleware function that uses user ID for rate limiting
func (rl *RateLimiter) RateLimitByUser() gin.HandlerFunc {
	go rl.cleanup() // Start cleanup goroutine

	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			// If no user ID, fall back to IP-based rate limiting
			key := c.ClientIP()
			limiter := rl.getLimiter(key)

			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"success": false,
					"message": "rate limit exceeded",
				})
				c.Abort()
				return
			}
		} else {
			// Use user ID for rate limiting
			key := string(userID)
			limiter := rl.getLimiter(key)

			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"success": false,
					"message": "rate limit exceeded",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
