resource "aws_route53_zone" "main" {
  name = var.domain_name
}

module "network" {
  source = "../../modules/network"

  project_name          = var.project_name
  environment           = var.environment
  vpc_cidr              = var.vpc_cidr
  availability_zones    = var.availability_zones
  enable_private_egress = var.enable_private_egress
}

module "security" {
  source = "../../modules/security"

  project_name   = var.project_name
  environment    = var.environment
  vpc_id         = module.network.vpc_id
  container_port = var.container_port
}

module "database" {
  source = "../../modules/database"

  project_name               = var.project_name
  environment                = var.environment
  private_subnet_ids         = module.network.private_subnet_ids
  kms_key_arn                = module.security.kms_key_arn
  db_instance_class          = var.db_instance_class
  db_name                    = var.db_name
  db_username                = var.db_username
  database_security_group_id = module.security.database_security_group_id
}

module "storage" {
  source = "../../modules/storage"

  project_name = var.project_name
  environment  = var.environment
  kms_key_arn  = module.security.kms_key_arn
}

module "application" {
  source = "../../modules/application"

  project_name           = var.project_name
  environment            = var.environment
  domain_name            = var.domain_name
  route53_zone_id        = aws_route53_zone.main.zone_id
  vpc_id                 = module.network.vpc_id
  public_subnet_ids      = module.network.public_subnet_ids
  private_subnet_ids     = module.network.private_subnet_ids
  ecs_task_cpu           = var.ecs_task_cpu
  ecs_task_memory        = var.ecs_task_memory
  ecs_desired_count      = var.ecs_desired_count
  container_port         = var.container_port
  image_tag              = var.image_tag
  frontend_image_tag     = var.frontend_image_tag
  frontend_desired_count = var.frontend_desired_count
  db_endpoint            = module.database.db_endpoint
  db_name                = var.db_name
  kms_key_arn            = module.security.kms_key_arn
  db_secret_arn          = module.database.db_secret_arn
  jwt_secret_arn         = module.security.jwt_secret_arn
  cors_allowed_origins   = "https://${var.domain_name},https://www.${var.domain_name}"
  s3_bucket_name         = module.storage.medical_images_bucket_name
  s3_bucket_arn          = module.storage.medical_images_bucket_arn
  audit_logs_bucket_name = module.storage.audit_logs_bucket_name
  alb_security_group_id  = module.security.alb_security_group_id
  ecs_security_group_id  = module.security.ecs_security_group_id
}

resource "aws_route53_record" "apex" {
  zone_id = aws_route53_zone.main.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = module.application.alb_dns_name
    zone_id                = module.application.alb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "www" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "www"
  type    = "A"

  alias {
    name                   = module.application.alb_dns_name
    zone_id                = module.application.alb_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "api"
  type    = "A"

  alias {
    name                   = module.application.alb_dns_name
    zone_id                = module.application.alb_zone_id
    evaluate_target_health = true
  }
}

module "observability" {
  source = "../../modules/observability"

  project_name           = var.project_name
  environment            = var.environment
  vpc_id                 = module.network.vpc_id
  audit_logs_bucket_name = module.storage.audit_logs_bucket_name
}
