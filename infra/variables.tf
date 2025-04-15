variable "app_port" {
  description = "アプリケーションのポート"
  default     = 3000
}

variable "app_image_tag" {
  description = "アプリケーションのイメージタグ"
  default     = "latest"
}

variable "db_username" {
  description = "Database master username"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "emotion_analysis"
}