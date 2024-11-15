package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger is a middleware that logs HTTP requests
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		entry := logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency":     latency,
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
			"user_agent":  userAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			msg := "Request processed"
			if statusCode >= 500 {
				entry.Error(msg)
			} else if statusCode >= 400 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
