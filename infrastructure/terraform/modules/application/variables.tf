variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name"
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

variable "db_secret_arn" {
  description = "Secrets Manager ARN for DB credentials"
  type        = string
}

variable "s3_bucket_arn" {
  description = "S3 bucket ARN for medical images"
  type        = string
}
