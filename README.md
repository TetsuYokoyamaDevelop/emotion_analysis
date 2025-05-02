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

Go 1.23

Gin Web Framework

GORM (ORM)

OpenAI API (Function Calling)

MySQL

Docker / Docker Compose

Terraform (AWS ECS/ECR/Fargate等)

GitHub Actions

### 使用技術選定の理由

Go 1.22
高速な処理性能とシンプルな構文に魅力を感じ、初めて本格的にGoを用いたプロダクトを開発しました。学習目的だけでなく、将来的に案件で活用できるよう、実践的なアプリケーションの構築を通して習得を目指しました。

Gin Web Framework
Go言語での開発において最も普及している軽量フレームワークであり、公式ドキュメントや記事も充実していることから採用しました。学習コストが低く、シンプルにAPIを構築できる点が大きなメリットです。

GORM
マイグレーションやリレーションの管理など、実用的なORM機能を備えており、スムーズにDB操作を行える点が魅力です。RailsやLaravelなど他言語でORMを扱っていた経験があり、親和性が高いと判断しました。

OpenAI API（Function Calling）
自然言語処理による感情分析を実装するためにOpenAI APIを採用しました。Function Calling機能を活用することで、AIの出力をJSON形式で安全に受け取り、サーバーサイドと正確に連携できる構成を実現しています。

MySQL
過去の開発経験から操作に慣れており、安定性と情報の多さから採用しました。

Docker / Docker Compose
環境差異をなくし、ローカル・本番問わず同一の構成で実行可能にするために導入しました。コンテナベースの構成により、開発・テスト・デプロイの手間を最小限に抑えています。

Terraform（AWS ECS/ECR/Fargateなど）
初学者ながらインフラのコード化に挑戦するため、Terraformを採用しました。Terraformの大部分はAIアシストを活用し、開発リソースをアプリケーション側に集中できるようにしています。ECS/Fargate構成はGoのWebアプリケーションと親和性が高く、スケーラブルな構成を意識しました。

GitHub Actions
CI/CDの自動化のために採用しました。mainブランチにマージすると自動的にECSへデプロイされる仕組みを構築し、継続的デリバリーを実現しています。

### ブラウザーで動作確認

/analyze にPOSTして感情分析結果を取得する

/users/registration でユーザー登録

/users/login でJWT発行

/history にGETして履歴を取得

### 注意

.env や GitHub Secretsで別途API KeyとDB設定の管理が必要

JWT_SECRET、PRIVATE_API_KEY、CUSTOM_OPENAI_KEYなど必須

目標：簡単に開発できて、すぐにデプロイできる感情分析APIプロジェクト