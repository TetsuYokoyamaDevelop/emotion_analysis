// internal/service/emotion_test.go
package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/TetsuYokoyamaDevelop/emotion_analysis.git/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupTestDB(t *testing.T) *gorm.DB {
	// 環境変数から接続情報を取得
	dbHost := getEnvOrDefault("TEST_DB_HOST", "localhost")
	dbPort := getEnvOrDefault("TEST_DB_PORT", "3306")
	dbUser := getEnvOrDefault("TEST_DB_USER", "root")
	dbPassword := getEnvOrDefault("TEST_DB_PASSWORD", "")
	dbName := getEnvOrDefault("TEST_DB_NAME", "emotion_analysis_test")

	// テスト用のデータベース接続文字列
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// データベース接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("データベース接続エラー: %v\nDSN: %s", err, dsn)
	}

	// テーブルをクリーンアップ
	err = db.Migrator().DropTable(&model.Message{}, &model.User{})
	assert.NoError(t, err)

	// マイグレーション実行
	err = db.AutoMigrate(&model.User{}, &model.Message{})
	assert.NoError(t, err)

	// テストユーザー作成
	user := model.User{
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	return db
}

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

func TestAnalyzeText(t *testing.T) {
	// テスト用のモックサーバー
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// モックレスポンス
		response := model.OpenAIResponse{
			Choices: []struct {
				Message struct {
					ToolCalls []struct {
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"message"`
			}{
				{
					Message: struct {
						ToolCalls []struct {
							Function struct {
								Name      string `json:"name"`
								Arguments string `json:"arguments"`
							} `json:"function"`
						} `json:"tool_calls"`
					}{
						ToolCalls: []struct {
							Function struct {
								Name      string `json:"name"`
								Arguments string `json:"arguments"`
							} `json:"function"`
						}{
							{
								Function: struct {
									Name      string `json:"name"`
									Arguments string `json:"arguments"`
								}{
									Name:      "analyze_sentiment",
									Arguments: `{"sentiment":"positive","sentimentScore":0.8,"explanation":"テスト説明","praise_or_advice":"テストアドバイス"}`,
								},
							},
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// テスト用のテンプレートファイル作成
	templateContent := `{
		"model": "gpt-4",
		"temperature": 0.7,
		"tools": [
			{
				"type": "function",
				"function": {
					"name": "analyze_sentiment",
					"description": "感情分析の結果を返す",
					"parameters": {
						"type": "object",
						"properties": {
							"sentiment": {
								"type": "string",
								"enum": ["positive", "neutral", "negative"],
								"description": "感情の種類"
							},
							"sentimentScore": {
								"type": "number",
								"description": "感情スコア（-1から1）"
							},
							"explanation": {
								"type": "string",
								"description": "詳細な説明"
							},
							"praise_or_advice": {
								"type": "string",
								"description": "褒め言葉またはアドバイス"
							}
						},
						"required": ["sentiment", "sentimentScore", "explanation", "praise_or_advice"]
					}
				}
			}
		]
	}`

	// テンプレートファイルのディレクトリを作成
	err := os.MkdirAll("internal/service", 0755)
	assert.NoError(t, err)

	// テンプレートファイルを作成
	templatePath := filepath.Join("internal/service", "prompt_template.json")
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(templatePath)

	tests := []struct {
		name           string
		text           string
		userEmail      string
		config         Config
		expectedResult model.SentimentResult
		expectError    bool
	}{
		{
			name:      "APIキー未設定",
			text:      "テストテキスト",
			userEmail: "test@example.com",
			config: Config{
				APIEndpoint: server.URL,
				APIKey:      "",
			},
			expectedResult: model.SentimentResult{
				Explanation: "APIキーが設定されていません",
			},
			expectError: false,
		},
		{
			name:      "正常系",
			text:      "テストテキスト",
			userEmail: "test@example.com",
			config: Config{
				APIEndpoint: server.URL,
				APIKey:      "test-api-key",
			},
			expectedResult: model.SentimentResult{
				Sentiment:      "positive",
				SentimentScore: 0.8,
				Explanation:    "テスト説明",
				PraiseOrAdvice: "テストアドバイス",
			},
			expectError: false,
		},
		{
			name:      "存在しないユーザー",
			text:      "テストテキスト",
			userEmail: "nonexistent@example.com",
			config: Config{
				APIEndpoint: server.URL,
				APIKey:      "test-api-key",
			},
			expectedResult: model.SentimentResult{
				Explanation: "データベース保存に失敗しました",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			result := AnalyzeText(tt.text, tt.userEmail, db, tt.config)

			if tt.expectError {
				assert.NotEmpty(t, result.Explanation)
				assert.Equal(t, tt.expectedResult.Explanation, result.Explanation)
			} else {
				assert.Equal(t, tt.expectedResult.Sentiment, result.Sentiment)
				assert.Equal(t, tt.expectedResult.SentimentScore, result.SentimentScore)
				assert.Equal(t, tt.expectedResult.Explanation, result.Explanation)
				assert.Equal(t, tt.expectedResult.PraiseOrAdvice, result.PraiseOrAdvice)
			}
		})
	}
}
