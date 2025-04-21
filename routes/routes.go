package routes

import (
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/analyze", handler.AnalyzeHandler)

	return r
}
