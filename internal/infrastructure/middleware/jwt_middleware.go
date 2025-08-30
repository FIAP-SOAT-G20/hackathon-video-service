package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/domain"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/core/port"
)

func JWTAuthMiddleware(jwtService port.JWTService) gin.HandlerFunc {
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

		if err := jwtService.ValidateToken(parts[1]); err != nil {
			_ = c.Error(domain.NewUnauthorizedError(domain.ErrInvalidToken))
			c.Abort()
			return
		}

		c.Next()
	}
}
