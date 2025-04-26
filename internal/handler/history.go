package handler

import (
	"fmt"
	"net/http"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HistoryHandler struct {
	DB *gorm.DB
}

// ここでPOST/analyzeを受ける
func (h AnalyzeHandler) GetHistory(c *gin.Context) {
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
	result, err := service.History(userEmailStr, h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history"})
		return
	}
	c.JSON(http.StatusOK, result)
}
