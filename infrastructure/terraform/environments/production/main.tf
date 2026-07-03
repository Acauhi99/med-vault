module "network" {
  source = "../../modules/network"

  project_name      = var.project_name
  environment       = var.environment
  vpc_cidr          = var.vpc_cidr
  availability_zones = var.availability_zones
}

module "security" {
  source = "../../modules/security"

  project_name = var.project_name
  environment  = var.environment
}

module "database" {
  source = "../../modules/database"

  project_name       = var.project_name
  environment        = var.environment
  vpc_id             = module.network.vpc_id
  private_subnet_ids = module.network.private_subnet_ids
  db_instance_class  = var.db_instance_class
  db_name            = var.db_name
  db_username        = var.db_username
  ecs_security_group_id = module.application.ecs_security_group_id
}

module "storage" {
  source = "../../modules/storage"

  project_name = var.project_name
  environment  = var.environment
}

module "application" {
  source = "../../modules/application"

  project_name       = var.project_name
  environment        = var.environment
  vpc_id             = module.network.vpc_id
  public_subnet_ids  = module.network.public_subnet_ids
  private_subnet_ids = module.network.private_subnet_ids
  ecs_task_cpu       = var.ecs_task_cpu
  ecs_task_memory    = var.ecs_task_memory
  ecs_desired_count  = var.ecs_desired_count
  container_port     = var.container_port
  db_secret_arn      = module.database.db_secret_arn
  s3_bucket_arn      = module.storage.medical_images_bucket_arn
}

module "observability" {
  source = "../../modules/observability"

  project_name = var.project_name
  environment  = var.environment
}
