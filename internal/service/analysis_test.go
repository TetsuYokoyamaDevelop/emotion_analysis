// internal/service/test_helper.go
package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockOpenAIServer はOpenAI APIのモックサーバーを作成する
func MockOpenAIServer(t *testing.T, sentiment string, score float64, explanation, advice string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// モックレスポンスを作成
		mockResponse := createMockOpenAIResponse(sentiment, score, explanation, advice)

		responseJSON, err := json.Marshal(mockResponse)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}))
}

// createMockOpenAIResponse はテスト用のOpenAIレスポンスを作成する
func createMockOpenAIResponse(sentiment string, score float64, explanation, advice string) map[string]interface{} {
	arguments := map[string]interface{}{
		"sentiment":        sentiment,
		"sentimentScore":   score,
		"explanation":      explanation,
		"praise_or_advice": advice,
	}

	argumentsJSON, _ := json.Marshal(arguments)

	return map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]interface{}{
					"tool_calls": []map[string]interface{}{
						{
							"function": map[string]interface{}{
								"arguments": string(argumentsJSON),
							},
						},
					},
				},
			},
		},
	}
}

// compareFloat64 は浮動小数点数の比較を行うヘルパー関数
func compareFloat64(t *testing.T, expected, actual float64, tolerance float64) bool {
	diff := expected - actual
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Errorf("Expected %f, got %f (difference: %f)", expected, actual, diff)
		return false
	}
	return true
}
