package main

import (
	"log"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/config"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/routes"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()
	db := config.InitDB()
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Message{}).Error; err != nil {
		log.Fatalf("Messageテーブルの削除に失敗しました: %v", err)
	}
	engine := routes.SetupRouter(db)
	engine.Run("0.0.0.0:3000")
}
