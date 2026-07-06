output "vpc_id" {
  description = "VPC ID"
  value       = module.network.vpc_id
}

output "route53_zone_name_servers" {
  description = "Name servers for the primary Route 53 hosted zone"
  value       = aws_route53_zone.main.name_servers
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS name"
  value       = module.application.alb_dns_name
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.application.ecs_cluster_name
}

output "ecs_service_name" {
  description = "ECS service name"
  value       = module.application.ecs_service_name
}

output "ecr_repository_url" {
  description = "ECR repository URL for backend Docker images"
  value       = module.application.ecr_repository_url
}

output "frontend_ecr_repository_url" {
  description = "ECR repository URL for frontend Docker images"
  value       = module.application.frontend_ecr_repository_url
}

output "frontend_ecs_service_name" {
  description = "Frontend ECS service name"
  value       = module.application.frontend_ecs_service_name
}

output "api_base_url" {
  description = "Public backend API base URL"
  value       = "https://api.${var.domain_name}"
}

output "frontend_base_url" {
  description = "Public frontend base URL"
  value       = "https://${var.domain_name}"
}

output "db_endpoint" {
  description = "RDS PostgreSQL endpoint"
  value       = module.database.db_endpoint
  sensitive   = true
}

output "s3_medical_images_bucket" {
  description = "S3 bucket for medical images"
  value       = module.storage.medical_images_bucket_name
}

output "cloudtrail_arn" {
  description = "CloudTrail ARN"
  value       = module.observability.cloudtrail_arn
}

output "vpc_flow_logs_group_name" {
  description = "VPC Flow Logs CloudWatch log group name"
  value       = module.observability.vpc_flow_logs_group_name
}
