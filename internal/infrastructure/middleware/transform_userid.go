package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TransformUserIDFromContextToQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized - user_id not found in context"})
			return
		}

		// Safe type assertion with check
		userIDUint64, ok := userID.(uint64)
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{"error": "Invalid user_id type - expected uint64"})
			return
		}

		query := c.Request.URL.Query()
		query.Set("user_id", fmt.Sprintf("%d", userIDUint64))
		c.Request.URL.RawQuery = query.Encode()

		c.Next()
	}
}
