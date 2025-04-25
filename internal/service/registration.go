package service

import (
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserRegist(db *gorm.DB, user model.User) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	user.Password = string(hashedPassword)

	if err := db.Save(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}
