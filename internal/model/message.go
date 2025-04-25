package model

import "time"

// Message はチャットメッセージを表す構造体
type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"-"`
	Role           string    `gorm:"size:20;not null;check:role IN ('user', 'assistant')" json:"role"`
	Sentiment      string    `gorm:"size:50" json:"sentiment"`          // 感情分類（例: "positive"）
	SentimentScore float64   `gorm:"type:float" json:"sentimentScore"`  // 感情スコア（例: 0.85）
	Explanation    string    `gorm:"type:text" json:"explanation"`      // 解説・補足
	PraiseOrAdvice string    `gorm:"type:text" json:"praise_or_advice"` // アドバイス or 賞賛文
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
