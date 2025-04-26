package main

import (
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/config"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/routes"
)

func main() {
	config.LoadEnv()
	db := config.InitDB()
	engine := routes.SetupRouter(db)
	engine.Run("0.0.0.0:3000")
}
