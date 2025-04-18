package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
)

func AnalyzeText(text string) model.SentimentResult {
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	if apiKey == "" {
		fmt.Println("APIキーが設定されてません")
		return model.SentimentResult{Explanation: "APIキーが設定されていません"}
	}

	// JSONテンプレート読み込み
	templatePath := "internal/service/prompt_template.json"
	data, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("テンプレート読み込み失敗:", err)
		return model.SentimentResult{Explanation: "テンプレートファイル読み込みに失敗しました"}
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		fmt.Println("テンプレートパース失敗:", err)
		return model.SentimentResult{Explanation: "テンプレートパースに失敗しました"}
	}

	// メッセージ追加
	payload["messages"] = []map[string]string{
		{
			"role":    "system",
			"content": "あなたは入力された文章の感情を分析し、感情の種類（positive/neutral/negative）、スコア（-1から1）、説明、そして感情に応じたメッセージを返してください。ポジティブな場合はユーザーを褒める言葉を、ネガティブな場合は前向きになれるアドバイスを、日本語で返してください。",
		},
		{
			"role":    "user",
			"content": text,
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("OpenAI APIエラー:", err)
		return model.SentimentResult{Explanation: "OpenAI API呼び出しに失敗しました"}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var aiResp model.OpenAIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		fmt.Println("JSONパース失敗:", err)
		return model.SentimentResult{Explanation: "レスポンスの解析に失敗しました"}
	}

	raw := aiResp.Choices[0].Message.ToolCalls[0].Function.Arguments

	var result model.SentimentResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		fmt.Println("Unmarshal失敗:", err)
		fmt.Println("中身:", raw)
		return model.SentimentResult{
			Explanation:    "出力の整形に失敗しました",
			Sentiment:      "unknown",
			SentimentScore: 0,
			PraiseOrAdvice: "",
		}
	}

	return result
}
