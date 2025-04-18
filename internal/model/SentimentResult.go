package model

type SentimentResult struct {
	Sentiment      string  `json:"sentiment"`
	SentimentScore float64 `json:"sentimentScore"`
	Explanation    string  `json:"explanation"`
	PraiseOrAdvice string  `json:"praise_or_advice"`
}
