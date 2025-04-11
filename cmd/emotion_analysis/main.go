package main

import "github.com/TetsuYokoyamaDevelop/emotion_analysis.git/routes"

func main() {
	engine := routes.SetupRouter()
	engine.Run("0.0.0.0:3000")
}
