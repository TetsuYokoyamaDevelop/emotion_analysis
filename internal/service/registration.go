package service

import (
	"net/http"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserRegist(user model.User, db *gorm.DB, c *gin.Context) (model.User, error) {
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return model.User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return model.User{}, err
	}
	user.Password = string(hashedPassword)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while registering user"})
		return model.User{}, err
	}
	return user, nil
}
