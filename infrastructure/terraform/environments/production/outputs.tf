output "vpc_id" {
  description = "VPC ID"
  value       = module.network.vpc_id
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS name"
  value       = module.application.alb_dns_name
}

output "ecr_repository_url" {
  description = "ECR repository URL for backend Docker images"
  value       = module.application.ecr_repository_url
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
