// internal/service/emotion_test.go
package service

import (
	"testing"
)

// 基本的なテストで動作確認
func TestBasic(t *testing.T) {
	t.Log("Test is running")
}

// APIキーのチェックのみのシンプルなテスト
func TestAnalyzeText_MissingAPIKey(t *testing.T) {
	config := Config{
		APIEndpoint: "http://example.com",
		APIKey:      "", // 空のAPIキー
	}

	// データベースはnilで渡す（このテストでは使用しない）
	result := AnalyzeText("テスト", "test@example.com", nil, config)

	if result.Explanation != "APIキーが設定されていません" {
		t.Errorf("Expected 'APIキーが設定されていません', got '%s'", result.Explanation)
	}
}
