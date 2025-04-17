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

func AnalyzeText(text string) model.Message {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("APIキーが設定されてません")
		return model.Message{Text: "APIキーが設定されていません"}
	}

	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo-0125",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "あなたは入力された文章の感情を分析し、感情の種類（positive/neutral/negative）、スコア（-1から1）、説明、そして感情に応じたメッセージを返してください。ポジティブな場合はユーザーを褒める言葉を、ネガティブな場合は前向きになれるアドバイスを、日本語で返してください。",
			},
			{
				"role":    "user",
				"content": text,
			},
		},
		"tools": []map[string]interface{}{
			{
				"type": "function",
				"function": map[string]interface{}{
					"name":        "analyze_sentiment_with_message",
					"description": "感情分析を行い、スコア・説明・適切なフィードバック（褒めorアドバイス）を返します。",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"sentiment": map[string]interface{}{
								"type": "string",
								"enum": []string{"positive", "neutral", "negative"},
							},
							"sentimentScore": map[string]interface{}{
								"type":        "number",
								"description": "-1から1のスコア",
							},
							"explanation": map[string]interface{}{
								"type":        "string",
								"description": "なぜこの感情とスコアになったか",
							},
							"praise_or_advice": map[string]interface{}{
								"type":        "string",
								"description": "ポジティブな場合はユーザーを褒める。ネガティブな場合は前向きになれるアドバイス。",
							},
						},
						"required": []string{"sentiment", "sentimentScore", "explanation", "praise_or_advice"},
					},
				},
			},
		},
		"tool_choice": map[string]interface{}{
			"type": "function",
			"function": map[string]string{
				"name": "analyze_sentiment_with_message",
			},
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
		return model.Message{Text: "OpenAI API呼び出しに失敗しました"}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var aiResp model.OpenAIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		fmt.Println("JSONパース失敗:", err)
		return model.Message{Text: "レスポンスの解析に失敗しました"}
	}

	var result model.Message
	if err := json.Unmarshal([]byte(aiResp.Choices[0].Message.FunctionCall.Arguments), &result); err != nil {
		fmt.Println("結果のUnmarshal失敗:", err)
		return model.Message{Text: "出力の整形に失敗しました"}
	}

	return result
}
