package service

import "github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"

// 受け取った文章を処理
func AnalyzeText(text string) model.AnalyzeResult {
	// ★ここは仮実装、今は固定値で返す
	// 本来はAPI叩いたり、テキスト解析する処理を入れる
	return model.AnalyzeResult{
		Sentiment: "positice",
		Score:     0.07,
	}
}
