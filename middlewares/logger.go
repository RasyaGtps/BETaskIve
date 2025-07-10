package middlewares

import (
	"time"
	"taskive/utils"

	"github.com/gin-gonic/gin"
)

var logger = utils.NewLogger()

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)

		logger.LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
} 