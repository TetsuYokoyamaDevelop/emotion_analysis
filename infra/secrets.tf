# データベース認証情報用のシークレット
resource "aws_secretsmanager_secret" "db_credentials" {
  name                    = "emotion-analysis/db-credentials"
  recovery_window_in_days = 0  # 開発中は削除を容易にするため（本番環境では7以上を推奨）
  
  tags = {
    Environment = "development"
    Application = "emotion-analysis"
  }
}

# シークレットの値
resource "aws_secretsmanager_secret_version" "db_credentials" {
  secret_id     = aws_secretsmanager_secret.db_credentials.id
  secret_string = jsonencode({
    username = var.db_username
    password = var.db_password
    host     = aws_db_instance.db.address
    port     = "3306"
    dbname   = var.db_name
  })
}