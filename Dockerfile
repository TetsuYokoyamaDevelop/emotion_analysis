# ベースイメージとして公式のGolangイメージを使用
FROM golang:1.23 AS builder

# 作業ディレクトリを設定
WORKDIR /app

# Go Modulesの使用を許可（必要な場合）
# ENV GO111MODULE=on

# ローカルのモジュールキャッシュを最適化（必要な場合）
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# アプリケーションをビルド
RUN go build -o myapp ./cmd/emotion_analysis

# 実行可能ファイルをデフォルトのコマンドとして設定
CMD ["./myapp"]

EXPOSE 3000