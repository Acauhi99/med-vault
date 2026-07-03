# Database module — placeholder for Phase 2
# Will contain: RDS PostgreSQL, subnet groups, parameter groups, security groups

locals {
  name_prefix = "${var.project_name}-${var.environment}"
}
