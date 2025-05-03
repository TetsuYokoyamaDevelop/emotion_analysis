package main

import (
	"log"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/config"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/routes"
)

func main() {
	config.LoadEnv()
	db := config.InitDB()
	engine := routes.SetupRouter(db)
	if err := engine.Run("0.0.0.0:3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
