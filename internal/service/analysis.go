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

	// JSONテンプレート読み込み
	templatePath := "internal/service/prompt_template.json"
	data, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("テンプレート読み込み失敗:", err)
		return model.Message{Text: "テンプレートファイル読み込みに失敗しました"}
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		fmt.Println("テンプレートパース失敗:", err)
		return model.Message{Text: "テンプレートパースに失敗しました"}
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
