package model

import "time"

// Message はチャットメッセージを表す構造体
type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"-"`
	Role           int       `gorm:"not null;check:role IN (1,2)" json:"role"` // 1=user, 2=assistant
	Sentiment      string    `gorm:"size:50" json:"sentiment"`
	SentimentScore float64   `gorm:"type:float" json:"sentimentScore"`
	Explanation    string    `gorm:"type:text" json:"explanation"`
	PraiseOrAdvice string    `gorm:"type:text" json:"praise_or_advice"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
