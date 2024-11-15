package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns the default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Total-Count",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// CORS middleware function
func CORS(config *CORSConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			c.Next()
			return
		}

		// Check if origin is allowed
		allowedOrigin := "*"
		for _, o := range config.AllowOrigins {
			if o == "*" || o == origin {
				allowedOrigin = o
				break
			}
			if strings.HasSuffix(o, "*") {
				prefix := strings.TrimSuffix(o, "*")
				if strings.HasPrefix(origin, prefix) {
					allowedOrigin = origin
					break
				}
			}
		}

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if allowedOrigin != "*" && len(config.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ","))
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
			if config.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", string(config.MaxAge))
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
