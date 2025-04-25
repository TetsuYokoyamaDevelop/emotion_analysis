package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 環境変数を取得（デフォルト値対応）
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 現在の環境を取得
func getEnvironment() string {
	return strings.ToLower(getEnv("APP_ENV", "development"))
}

// LoadEnv 環境変数を.envファイルから読み込む
func LoadEnv() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Printf("警告: .envファイルが見つからないか、読み込めませんでした")
		} else {
			log.Printf(".envファイルを読み込みました")
		}
	}
}

func InitDB() *gorm.DB {
	// 環境変数から取得（デフォルト値あり）
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "password")
	dbName := getEnv("DB_NAME", "emotion_analysis_dev")

	env := getEnvironment()
	log.Printf("Environment: %s", env)
	log.Printf("Connecting to database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempt %d to connect to the database", i+1)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Printf("Successfully connected to the database")

			// AutoMigrateをここで実行（構造体の定義を反映）
			err = db.AutoMigrate(&model.User{}, &model.Message{})
			if err != nil {
				log.Fatalf("AutoMigrate failed: %v", err)
			}
			log.Printf("AutoMigrate completed successfully")

			return db
		}
		log.Printf("Failed to connect to database: %v", err)
		log.Printf("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	return nil
}
