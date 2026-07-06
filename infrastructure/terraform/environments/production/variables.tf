variable "aws_region" {
  description = "AWS region for all resources"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name used for resource naming and tagging"
  type        = string
  default     = "medvault"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones for subnet distribution"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

variable "enable_private_egress" {
  description = "Create NAT gateway and private default route for ECS outbound access"
  type        = bool
  default     = false
}

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t4g.micro"
}

variable "db_name" {
  description = "PostgreSQL database name"
  type        = string
  default     = "medvault"
}

variable "db_username" {
  description = "PostgreSQL master username"
  type        = string
  default     = "medvault"
}

variable "ecs_task_cpu" {
  description = "ECS Fargate task CPU units"
  type        = number
  default     = 256
}

variable "ecs_task_memory" {
  description = "ECS Fargate task memory (MiB)"
  type        = number
  default     = 512
}

variable "ecs_desired_count" {
  description = "Number of ECS tasks to run"
  type        = number
  default     = 0
}

variable "container_port" {
  description = "Container port for the Go backend"
  type        = number
  default     = 8080
}

variable "image_tag" {
  description = "Backend Docker image tag"
  type        = string
  default     = "bootstrap"
}
