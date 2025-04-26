package routes

import (
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/handler"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	userHandler := handler.UserHandler{DB: db}
	loginHandler := handler.LoginHandler{DB: db}
	analyzeHandler := handler.AnalyzeHandler{DB: db}
	// ログイン (ミドルウェアなし)
	r.POST("/users/login", loginHandler.UserLoginHandler)
	// ユーザー登録 (APIキー認証)
	r.POST("/users/registration", middleware.AuthMiddleware(), userHandler.UserRegistHandler)
	// 感情分析 (JWTトークン認証)
	r.POST("/analyze", middleware.TokenMiddleware(), analyzeHandler.Analyze)

	r.GET("/history", middleware.TokenMiddleware(), analyzeHandler.GetHistory)

	return r
}
