output "alb_dns_name" {
  description = "ALBのDNS名"
  value       = aws_lb.ecs_alb.dns_name
}

output "ecr_repository_url" {
  description = "ECRリポジトリのURL"
  value       = aws_ecr_repository.app.repository_url
}