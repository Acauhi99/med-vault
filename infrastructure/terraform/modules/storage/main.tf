# Storage module — placeholder for Phase 2
# Will contain: S3 buckets, lifecycle rules, bucket policies

locals {
  name_prefix = "${var.project_name}-${var.environment}"
}
