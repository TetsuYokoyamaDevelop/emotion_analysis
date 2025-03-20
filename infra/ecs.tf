resource "aws_ecs_cluster" "app_cluster" {
  name = "emotion_analysis-cluster"
}

resource "aws_ecs_task_definition" "app_task" {
  family                   = "emotion_analysis"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  cpu                      = "256"
  memory                   = "512"

  container_definitions = jsonencode([
    {
      name      = "app"
      image     = "${aws_ecr_repository.app.repository_url}:latest"
      essential = true
      portMappings = [
        {
          containerPort = 3000  # 8080から3000に変更
          hostPort      = 3000  # 8080から3000に変更
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/emotion_analysis"
          "awslogs-region"        = "ap-northeast-1"
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

resource "aws_ecs_service" "app_service" {
  name            = "emotion_analysis-service"
  cluster         = aws_ecs_cluster.app_cluster.id
  task_definition = aws_ecs_task_definition.app_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.public_1.id, aws_subnet.public_2.id]
    security_groups  = [aws_security_group.app_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.ecs_tg.arn
    container_name   = "app"
    container_port   = 3000  # 8080から3000に変更
  }

  # この行を追加して、サービスがリスナーとターゲットグループに依存していることを明示する
  depends_on = [aws_lb_listener.ecs_http, aws_lb_target_group.ecs_tg]
}

# CloudWatch Logsグループの追加
resource "aws_cloudwatch_log_group" "ecs_logs" {
  name              = "/ecs/emotion_analysis"
  retention_in_days = 30
}