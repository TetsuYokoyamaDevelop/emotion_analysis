package service

import (
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"gorm.io/gorm"
)

func History(userEmail string, db *gorm.DB) ([]model.Message, error) {
	// まずUserテーブルからIDを取得
	var user model.User
	if err := db.Where("email = ?", userEmail).First(&user).Error; err != nil {
		return nil, err
	}

	// そのユーザーIDに紐づくメッセージを取得（最新100件）
	var history []model.Message
	if err := db.Where("user_id = ?", user.ID).
		Order("created_at desc").
		Limit(100).
		Find(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}
