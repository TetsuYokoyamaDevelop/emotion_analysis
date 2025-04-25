package routes

import (
	"net/http"
	"os"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	userHandler := handler.UserHandler{DB: db}
	loginHandler := handler.LoginHandler{DB: db}
	// 保護された感情分析API
	r.POST("/analyze", AuthMiddleware(), handler.AnalyzeHandler)
	r.POST("/users/registration", AuthMiddleware(), userHandler.UserRegistHandler)
	r.POST("/users/login", AuthMiddleware(), loginHandler.UserLoginHandler)

	return r
}
