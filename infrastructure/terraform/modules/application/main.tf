# Application module — placeholder for Phase 2
# Will contain: ECS cluster, service, ALB, task definition, security groups

locals {
  name_prefix = "${var.project_name}-${var.environment}"
}
