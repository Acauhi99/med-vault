locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "kms" {
  #checkov:skip=CKV_AWS_356:KMS key policies require Resource="*"; access is constrained by principal.
  #checkov:skip=CKV_AWS_111:KMS key administration is constrained to account root principal.
  #checkov:skip=CKV_AWS_109:KMS key administration is constrained to account root principal.
  statement {
    sid     = "EnableAccountAdministration"
    actions = ["kms:*"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }

    resources = ["*"]
  }
}

resource "aws_kms_key" "main" {
  description             = "MedVault ${var.environment} encryption key"
  deletion_window_in_days = 7
  enable_key_rotation     = true
  policy                  = data.aws_iam_policy_document.kms.json

  tags = {
    Name = "${local.name_prefix}-kms"
  }
}

resource "aws_kms_alias" "main" {
  name          = "alias/${local.name_prefix}"
  target_key_id = aws_kms_key.main.key_id
}

resource "aws_security_group" "alb" {
  #checkov:skip=CKV2_AWS_5:Attached to ALB in application module; Checkov cannot infer cross-module attachment.
  name        = "${local.name_prefix}-alb-sg"
  description = "Allow public HTTPS access to ALB"
  vpc_id      = var.vpc_id

  #checkov:skip=CKV_AWS_260:Public HTTPS access is required for the documented ALB ingress path.
  ingress {
    description = "HTTPS from internet"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${local.name_prefix}-alb-sg"
  }
}

resource "aws_security_group" "ecs" {
  #checkov:skip=CKV2_AWS_5:Attached to ECS service in application module; Checkov cannot infer cross-module attachment.
  name        = "${local.name_prefix}-ecs-sg"
  description = "Allow ALB traffic to ECS tasks"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Backend traffic from ALB"
    from_port       = var.container_port
    to_port         = var.container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  tags = {
    Name = "${local.name_prefix}-ecs-sg"
  }
}

resource "aws_security_group" "database" {
  #checkov:skip=CKV2_AWS_5:Attached to RDS in database module; Checkov cannot infer cross-module attachment.
  name        = "${local.name_prefix}-database-sg"
  description = "Allow PostgreSQL access from ECS tasks"
  vpc_id      = var.vpc_id

  ingress {
    description     = "PostgreSQL from ECS"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  tags = {
    Name = "${local.name_prefix}-database-sg"
  }
}

resource "aws_security_group_rule" "alb_to_ecs" {
  type                     = "egress"
  description              = "ALB to ECS backend"
  from_port                = var.container_port
  to_port                  = var.container_port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.alb.id
  source_security_group_id = aws_security_group.ecs.id
}

resource "aws_security_group_rule" "ecs_to_database" {
  type                     = "egress"
  description              = "ECS to PostgreSQL"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  security_group_id        = aws_security_group.ecs.id
  source_security_group_id = aws_security_group.database.id
}

resource "aws_security_group_rule" "ecs_to_https" {
  type              = "egress"
  description       = "ECS to AWS APIs over HTTPS"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.ecs.id
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_secretsmanager_secret" "jwt_signing_key" {
  #checkov:skip=CKV2_AWS_57:JWT rotation needs application support; documented as future secrets-rotation work.
  name                    = "${local.name_prefix}/jwt-signing-key"
  description             = "JWT signing key for MedVault backend"
  kms_key_id              = aws_kms_key.main.arn
  recovery_window_in_days = 7

  tags = {
    Name = "${local.name_prefix}-jwt-signing-key"
  }
}
