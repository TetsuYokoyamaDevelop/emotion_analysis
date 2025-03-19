# Emotion Analysis
このプロジェクトは、Go言語で作成されたシンプルなアプリケーションです。
Dockerを使用して簡単にビルドおよび実行できます。

## ファイル構成
```
├── Dockerfile          # Golang公式イメージを使用したDocker設定
├── docker-compose.yml  # コンテナの構成ファイル
├── go.mod              # Go Modulesの設定ファイル
├── main.go             # メインのGoプログラム
├── .gitignore          # gitから除外したいものの設定
├── README.md           # 本ファイル
```

## セットアップ
### 1. Dockerを使用した実行
#### 1.1 Docker Compose を使用する場合
以下のコマンドでコンテナをビルドして実行できます。
```
docker-compose up --build
```
#### 1.2 Docker を単体で使用する場合
Dockerイメージをビルド
```
docker build -t emotion_analysis .
```
コンテナを実行
```
docker run --rm myapp
```

動作確認
実行後、以下のように Hello, World! が表示されることを確認してください。
```
Hello, World!
```