package middleware

import (
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		requestID := c.GetString("request_id")

		c.Next()

		// Check if response came from cache
		_, exists := c.Get("cache_miss")

		log.InfoContext(c.Request.Context(),
			"request completed",
			"method", c.Request.Method,
			"path", path,
			"raw_query", raw,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"client_ip", c.ClientIP(),
			"request_id", requestID,
			"from_cache", !exists,
			"errors", c.Errors.String(),
		)
	}
}
