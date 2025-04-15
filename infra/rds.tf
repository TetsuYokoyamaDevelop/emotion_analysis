# RDSセキュリティグループ
resource "aws_security_group" "db_sg" {
  name        = "database-sg"
  description = "Security group for database"
  vpc_id      = aws_vpc.main.id

  # ECSタスクからのアクセスのみ許可
  ingress {
    from_port       = 3306  # MySQLの場合
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.app_sg.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# RDSサブネットグループ
resource "aws_db_subnet_group" "db_subnet_group" {
  name       = "db-subnet-group"
  subnet_ids = [aws_subnet.public_1.id, aws_subnet.public_2.id]
}

# RDSインスタンス
resource "aws_db_instance" "db" {
  identifier             = "emotion-analysis-db"
  engine                 = "mysql"  # または "postgres"
  engine_version         = "8.0"    # 適切なバージョンに変更
  instance_class         = "db.t3.micro"
  allocated_storage      = 20
  storage_type           = "gp2"
  username               = var.db_username
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.db_subnet_group.name
  vpc_security_group_ids = [aws_security_group.db_sg.id]
  publicly_accessible    = false
  skip_final_snapshot    = true
  
  # DBの名前
  db_name                = var.db_name
}