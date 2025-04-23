package routes

import (
	"net/http"
	"os"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/handler"
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

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 保護された感情分析API
	r.POST("/analyze", AuthMiddleware(), handler.AnalyzeHandler)

	return r
}
