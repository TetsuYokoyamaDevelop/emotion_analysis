provider "aws" {
  region = "ap-northeast-1"  # 東京リージョン（適宜変更）
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.80.0"
    }
  }
  backend "s3" {
    bucket         = "emotion-analysis-backend"
    key            = "ecs/terraform.tfstate"
    region         = "ap-northeast-1"
    encrypt        = true
  }
}
