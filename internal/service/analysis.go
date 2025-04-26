package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"gorm.io/gorm"
)

const (
	ASSISTANT_ROLE = 2
	USER_ROLE      = 1
)

func AnalyzeText(text string, userEmail string, db *gorm.DB) model.SentimentResult {
	apiKey := os.Getenv("CUSTOM_OPENAI_KEY")
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
			"role": "system",
			"content": `あなたは、入力された日本語の文章からその感情を深く分析するAIカウンセラーです。

以下の4項目をJSON形式で返してください：

1. sentiment（感情の種類）：positive / neutral / negative
2. sentimentScore（感情スコア）：-1から1の実数値（小数点2桁まで）
3. explanation（詳細な説明）：1000文字以上で、文脈・語彙・心理状態を深く考察。
4. praise_or_advice（フィードバック）：1000文字以上の具体的な褒め言葉または前向きなアドバイス。

出力はエスケープされていないJSONオブジェクトとして返してください。`,
		},
		{
			"role":    "user",
			"content": text,
		},
	}

	payload["max_tokens"] = 1600 // 長文対応のためちょい増やし

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

	// Unquote（もしエスケープされてたら）
	unquoted, err := strconv.Unquote(raw)
	if err != nil {
		// Unquote しなくてよかった場合はそのまま raw を使う
		unquoted = raw
	}

	var result model.SentimentResult
	if err := json.Unmarshal([]byte(unquoted), &result); err != nil {
		fmt.Println("Unmarshal失敗:", err)
		fmt.Println("中身:", unquoted)
		return model.SentimentResult{
			Explanation:    "出力の整形に失敗しました",
			Sentiment:      "unknown",
			SentimentScore: 0,
			PraiseOrAdvice: "",
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Where("email = ?", userEmail).First(&user).Error; err != nil {
			fmt.Println("ユーザーIDが見つかりません:", err)
			return err
		}

		userMsg := model.Message{
			UserID:         user.ID,
			Role:           USER_ROLE,
			Sentiment:      result.Sentiment,
			SentimentScore: result.SentimentScore,
			Explanation:    text,
			PraiseOrAdvice: "",
		}

		if err := tx.Select("UserID", "Role", "Sentiment", "SentimentScore", "Explanation", "PraiseOrAdvice").
			Create(&userMsg).Error; err != nil {
			fmt.Println("ユーザーメッセージ保存失敗:", err)
			return err
		}

		aiMsg := model.Message{
			UserID:         user.ID,
			Role:           ASSISTANT_ROLE,
			Sentiment:      result.Sentiment,
			SentimentScore: result.SentimentScore,
			Explanation:    result.Explanation,
			PraiseOrAdvice: result.PraiseOrAdvice,
		}

		if err := tx.Select("UserID", "Role", "Sentiment", "SentimentScore", "Explanation", "PraiseOrAdvice").
			Create(&aiMsg).Error; err != nil {
			fmt.Println("AIメッセージ保存失敗:", err)
			return err
		}

		// 全部成功したらnil返す（→コミット）
		return nil
	})

	if err != nil {
		fmt.Println("トランザクション失敗:", err)
		return model.SentimentResult{Explanation: "データベース保存に失敗しました"}
	}

	return result
}
