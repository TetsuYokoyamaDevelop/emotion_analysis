package main

import (
	"fmt"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/config"
	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/routes"
)

func main() {
	db := config.InitDB()
	fmt.Println(db)
	engine := routes.SetupRouter()
	engine.Run("0.0.0.0:3000")
}
