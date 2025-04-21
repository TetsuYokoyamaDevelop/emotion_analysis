package model

import "time"

// Message はチャットメッセージを表す構造体
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Role      string    `gorm:"size:20;not null;check:role IN ('user', 'assistant')" json:"role"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	Score     float64   `gorm:"type:float" json:"score"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
