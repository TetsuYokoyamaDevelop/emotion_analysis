package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// データベース接続の初期化
func initDB() *gorm.DB {
	// 環境変数からDB接続情報を取得
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// マイグレーション
	db.AutoMigrate(&model.User{}, &model.Message{})

	return db
}

func main() {
	db := initDB()
	fmt.Println(db)
	engine := routes.SetupRouter()
	engine.Run("0.0.0.0:3000")
}
