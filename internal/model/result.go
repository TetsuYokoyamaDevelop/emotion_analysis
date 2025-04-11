package model

// 分析結果の型を書く
type AnalyzeRequest struct {
	Text string `json:"text" binding:"required"`
}

type AnalyzeResult struct {
	Sentiment string  `json:"sentinent"`
	Score     float64 `json:"score"`
}
