package handler

import (
	"net/http"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

// ここでPOST/analyzeを受ける
func (h *UserHandler) UserRegistHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := service.UserRegist(user, h.DB, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusOK, result)
}
