package middleware

import (
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.ErrorContext(c.Request.Context(),
					"panic recovered",
					"error", err,
					"request_id", c.GetString("request_id"),
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
