# Emotion Analysis

Go 言語で構築した感情分析APIサービスです。
OpenAI APIを利用し、入力されたテキストの感情分析結果をJSONで返信します。
Docker / ECS / GitHub Actions / Terraform を利用した自動デプロイパイプラインの構成になっています。

ファイル構成
```
.
├── cmd/
│   └── emotion_analysis/    # エントリーポイント（main.go）
├── config/                      # DB初期化や環境変数読込み
├── infra/                       # Terraform設定（AWS環境構築）
├── internal/
│   ├── handler/              # HTTPハンドラ
│   ├── middleware/           #認証ミドルウェア
│   ├── model/                # データベースモデル
│   └── service/              # ビジネスロジック
├── routes/                      # APIルーティング設定
├── .github/workflows/deploy.yml  # GitHub Actions デプロイワークフロー
├── ecs-task-def-template.json    # ECSタスク定義テンプレート
├── docker-compose.yml            # Docker構成
├── Dockerfile                    # Dockerイメージ構築
├── go.mod                        # Go Modules設定
├── README.md                     # 本ファイル
```
## セットアップ

### 1. Dockerを使用した実行

#### 1.1 Docker Compose の場合
```
docker-compose up --build
```

#### 1.2 Docker 単体の場合

Dockerイメージをビルド
```
docker build -t emotion_analysis .
```
コンテナを実行
```
docker run --rm emotion_analysis
```
### 2. 自動デプロイ

GitHubで main ブランチにマージされると

GitHub Actions (.github/workflows/deploy.yml) が動作

DockerイメージをECRにpush

ECS (Fargate)で新しいタスクを起動

タスク定義は ecs-task-def-template.jsonから作成

### キー技術

Go 1.22

Gin Web Framework

GORM (ORM)

OpenAI API (Function Calling)

MySQL

Docker / Docker Compose

Terraform (AWS ECS/ECR/Fargate等)

GitHub Actions

### ブラウザーで動作確認

/analyze にPOSTして感情分析結果を取得する

/users/registration でユーザー登録

/users/login でJWT発行

/history にGETして履歴を取得

### 注意

.env や GitHub Secretsで別途API KeyとDB設定の管理が必要

JWT_SECRET、PRIVATE_API_KEY、CUSTOM_OPENAI_KEYなど必須

目標：簡単に開発できて、すぐにデプロイできる感情分析APIプロジェクト