package middleware

import (
	"github.com/gin-gonic/gin"
)

func Next(next gin.HandlerFunc) gin.HandlerFunc {
	return next
}
