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

		query := c.Request.URL.Query()
		query.Set("user_id", fmt.Sprintf("%d", userID.(uint64)))
		c.Request.URL.RawQuery = query.Encode()

		c.Next()
	}
}
