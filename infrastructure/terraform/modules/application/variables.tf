variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "domain_name" {
  description = "Primary DNS name for the ALB certificate"
  type        = string
}

variable "route53_zone_id" {
  description = "Route 53 hosted zone ID for DNS validation"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs for ALB"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for ECS"
  type        = list(string)
}

variable "ecs_task_cpu" {
  description = "ECS task CPU units"
  type        = number
}

variable "ecs_task_memory" {
  description = "ECS task memory (MiB)"
  type        = number
}

variable "ecs_desired_count" {
  description = "Number of ECS tasks"
  type        = number
}

variable "container_port" {
  description = "Container port"
  type        = number
}

variable "image_tag" {
  description = "Backend Docker image tag"
  type        = string
  default     = "bootstrap"
}

variable "db_endpoint" {
  description = "RDS endpoint hostname"
  type        = string
}

variable "db_name" {
  description = "Database name"
  type        = string
}

variable "kms_key_arn" {
  description = "KMS key ARN for encryption"
  type        = string
}

variable "db_secret_arn" {
  description = "Secrets Manager ARN for DB credentials"
  type        = string
}

variable "jwt_secret_arn" {
  description = "JWT signing key secret ARN"
  type        = string
}

variable "s3_bucket_name" {
  description = "S3 bucket name for medical images"
  type        = string
}

variable "s3_bucket_arn" {
  description = "S3 bucket ARN for medical images"
  type        = string
}

variable "audit_logs_bucket_name" {
  description = "S3 bucket name for ALB access logs"
  type        = string
}

variable "alb_security_group_id" {
  description = "ALB security group ID"
  type        = string
}

variable "ecs_security_group_id" {
  description = "ECS security group ID"
  type        = string
}
