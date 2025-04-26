package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// APIキー認証ミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientKey := c.GetHeader("X-API-KEY")
		expectedKey := os.Getenv("PRIVATE_API_KEY")

		if expectedKey == "" || clientKey != expectedKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized – go away.",
			})
			return
		}
		c.Next()
	}
}
