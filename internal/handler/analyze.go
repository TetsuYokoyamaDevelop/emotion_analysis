package handler

import (
	"fmt"
	"net/http"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyzeHandler struct {
	DB *gorm.DB
}

// ここでPOST/analyzeを受ける
func (h AnalyzeHandler) Analyze(c *gin.Context) {
	var input struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userEmail, exists := c.Get("userEmail")
	if !exists {
		fmt.Println("ユーザーのメールが見つかりません")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	userEmailStr, ok := userEmail.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user email"})
		return
	}
	result := service.AnalyzeText(input.Text, userEmailStr, h.DB)
	c.JSON(http.StatusOK, result)
}
