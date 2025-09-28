package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
)

func JWTAuth(jwtService port.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(domain.NewUnauthorizedError(domain.ErrMissingAuthHeader))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			_ = c.Error(domain.NewUnauthorizedError(domain.ErrInvalidAuthHeader))
			c.Abort()
			return
		}

		if parts[1] == "" {
			_ = c.Error(domain.NewUnauthorizedError(domain.ErrInvalidToken))
			c.Abort()
			return
		}

		userID, err := jwtService.ExtractUserIDFromToken(parts[1])
		if err != nil {
			_ = c.Error(domain.NewUnauthorizedError(domain.ErrInvalidToken))
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
