package model

import (
	"time"
)

// User はユーザー情報を表す構造体
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"size:100;not null;unique"`
	Email        string    `gorm:"size:255;not null;unique"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Messages     []Message `gorm:"foreignKey:UserID"`
}

// Message はチャットメッセージを表す構造体
type Message struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             uint      `gorm:"not null;index"`
	User               User      `gorm:"foreignKey:UserID"`
	Role               string    `gorm:"size:20;not null;check:role IN ('user', 'assistant')"`
	Content            string    `gorm:"type:text;not null"`
	SentimentScore     float64   `gorm:"type:float"`
	SentimentMagnitude float64   `gorm:"type:float"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP;index"`
}
