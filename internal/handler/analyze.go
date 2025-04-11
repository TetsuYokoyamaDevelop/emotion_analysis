package handler

import (
	"net/http"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/service"
	"github.com/gin-gonic/gin"
)

// ここでPOST/analyzeを受ける
func AnalyzeHandler(c *gin.Context) {
	var input model.AnalyzeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := service.AnalyzeText(input.Text)
	c.JSON(http.StatusOK, result)
}
