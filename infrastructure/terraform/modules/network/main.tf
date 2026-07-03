# Network module — placeholder for Phase 2
# Will contain: VPC, public/private subnets, NAT gateway, route tables, IGW

locals {
  name_prefix = "${var.project_name}-${var.environment}"
}
