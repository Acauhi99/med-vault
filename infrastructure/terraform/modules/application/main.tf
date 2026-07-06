locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

data "aws_iam_policy_document" "ecs_tasks_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "task_execution_secrets" {
  statement {
    actions = [
      "secretsmanager:GetSecretValue",
    ]

    resources = [
      var.db_secret_arn,
      var.jwt_secret_arn,
    ]
  }
}

data "aws_iam_policy_document" "task_s3_access" {
  statement {
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]

    resources = ["${var.s3_bucket_arn}/*"]
  }

  statement {
    actions   = ["s3:ListBucket"]
    resources = [var.s3_bucket_arn]
  }
}

resource "aws_acm_certificate" "https" {
  domain_name = var.domain_name
  subject_alternative_names = [
    "www.${var.domain_name}",
    "api.${var.domain_name}",
  ]
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_ecr_repository" "frontend" {
  name                 = "${var.project_name}/frontend"
  image_tag_mutability = "IMMUTABLE"

  encryption_configuration {
    encryption_type = "KMS"
    kms_key         = var.kms_key_arn
  }

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${local.name_prefix}-frontend-ecr"
  }
}

resource "aws_cloudwatch_log_group" "frontend" {
  #checkov:skip=CKV_AWS_338:Project security doc sets CloudWatch application log retention to 90 days; audit records are retained in S3 for 6 years.
  name              = "/ecs/${local.name_prefix}-frontend"
  retention_in_days = 90

  tags = {
    Name = "${local.name_prefix}-frontend-logs"
  }
}

resource "aws_route53_record" "https_validation" {
  for_each = {
    for dvo in aws_acm_certificate.https.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id = var.route53_zone_id
  name    = each.value.name
  type    = each.value.type
  ttl     = 60
  records = [each.value.record]
}

resource "aws_acm_certificate_validation" "https" {
  certificate_arn         = aws_acm_certificate.https.arn
  validation_record_fqdns = [for record in aws_route53_record.https_validation : record.fqdn]
}

resource "aws_ecr_repository" "backend" {
  name                 = "${var.project_name}/backend"
  image_tag_mutability = "IMMUTABLE"

  encryption_configuration {
    encryption_type = "KMS"
    kms_key         = var.kms_key_arn
  }

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${local.name_prefix}-backend-ecr"
  }
}

resource "aws_cloudwatch_log_group" "backend" {
  #checkov:skip=CKV_AWS_338:Project security doc sets CloudWatch application log retention to 90 days; audit records are retained in S3 for 6 years.
  name              = "/ecs/${local.name_prefix}-backend"
  retention_in_days = 90

  tags = {
    Name = "${local.name_prefix}-backend-logs"
  }
}

resource "aws_iam_role" "task_execution" {
  name               = "${local.name_prefix}-ecs-task-execution"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_role.json
}

resource "aws_iam_role_policy_attachment" "task_execution" {
  role       = aws_iam_role.task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "task_execution_secrets" {
  name   = "${local.name_prefix}-ecs-task-secrets"
  role   = aws_iam_role.task_execution.id
  policy = data.aws_iam_policy_document.task_execution_secrets.json
}

resource "aws_iam_role" "task" {
  name               = "${local.name_prefix}-ecs-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume_role.json
}

resource "aws_iam_role_policy" "task_s3_access" {
  name   = "${local.name_prefix}-ecs-task-s3"
  role   = aws_iam_role.task.id
  policy = data.aws_iam_policy_document.task_s3_access.json
}

resource "aws_ecs_cluster" "main" {
  name = "${local.name_prefix}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "${local.name_prefix}-cluster"
  }
}

resource "aws_lb" "main" {
  #checkov:skip=CKV2_AWS_76:WAF has AWSManagedRulesKnownBadInputsRuleSet; Checkov does not infer it through this module association.
  name               = "${local.name_prefix}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_security_group_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = true
  drop_invalid_header_fields = true

  access_logs {
    bucket  = var.audit_logs_bucket_name
    prefix  = "alb"
    enabled = true
  }

  tags = {
    Name = "${local.name_prefix}-alb"
  }
}

resource "aws_lb_target_group" "backend" {
  #checkov:skip=CKV_AWS_378:ALB to ECS uses HTTP inside private VPC per documented architecture.
  name        = "${local.name_prefix}-backend-tg"
  port        = var.container_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = var.vpc_id

  health_check {
    enabled             = true
    path                = "/health"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
  }

  tags = {
    Name = "${local.name_prefix}-backend-tg"
  }
}

resource "aws_lb_target_group" "frontend" {
  name        = "${local.name_prefix}-frontend-tg"
  port        = var.container_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = var.vpc_id

  health_check {
    enabled             = true
    path                = "/"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
  }

  tags = {
    Name = "${local.name_prefix}-frontend-tg"
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.main.arn
  port              = 443
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate_validation.https.certificate_arn

  ssl_policy = "ELBSecurityPolicy-TLS13-1-2-2021-06"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.frontend.arn
  }
}

resource "aws_lb_listener_rule" "backend" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.backend.arn
  }

  condition {
    host_header {
      values = ["api.${var.domain_name}"]
    }
  }
}

resource "aws_ecs_task_definition" "backend" {
  family                   = "${local.name_prefix}-backend"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = tostring(var.ecs_task_cpu)
  memory                   = tostring(var.ecs_task_memory)
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.task.arn

  container_definitions = jsonencode([
    {
      name                   = "backend"
      image                  = "${aws_ecr_repository.backend.repository_url}:${var.image_tag}"
      essential              = true
      readonlyRootFilesystem = true

      portMappings = [
        {
          containerPort = var.container_port
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "DB_HOST"
          value = var.db_endpoint
        },
        {
          name  = "DB_NAME"
          value = var.db_name
        },
        {
          name  = "DB_PORT"
          value = "5432"
        },
        {
          name  = "S3_BUCKET"
          value = var.s3_bucket_name
        },
        {
          name  = "AWS_REGION"
          value = data.aws_region.current.name
        },
        {
          name  = "CORS_ALLOWED_ORIGINS"
          value = var.cors_allowed_origins
        },
      ]

      secrets = [
        {
          name      = "DB_USERNAME"
          valueFrom = "${var.db_secret_arn}:username::"
        },
        {
          name      = "DB_PASSWORD"
          valueFrom = "${var.db_secret_arn}:password::"
        },
        {
          name      = "JWT_SECRET"
          valueFrom = var.jwt_secret_arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.backend.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "backend"
        }
      }
    }
  ])

  tags = {
    Name = "${local.name_prefix}-backend-task"
  }
}

resource "aws_ecs_task_definition" "frontend" {
  family                   = "${local.name_prefix}-frontend"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = tostring(var.ecs_task_cpu)
  memory                   = tostring(var.ecs_task_memory)
  execution_role_arn       = aws_iam_role.task_execution.arn

  container_definitions = jsonencode([
    {
      name                   = "frontend"
      image                  = "${aws_ecr_repository.frontend.repository_url}:${var.frontend_image_tag}"
      essential              = true
      readonlyRootFilesystem = false

      portMappings = [
        {
          containerPort = var.container_port
          protocol      = "tcp"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.frontend.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "frontend"
        }
      }
    }
  ])

  tags = {
    Name = "${local.name_prefix}-frontend-task"
  }
}

data "aws_region" "current" {}

resource "aws_ecs_service" "backend" {
  name            = "${local.name_prefix}-backend"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.backend.arn
  desired_count   = var.ecs_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.ecs_security_group_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.backend.arn
    container_name   = "backend"
    container_port   = var.container_port
  }

  depends_on = [aws_lb_listener.https]

  tags = {
    Name = "${local.name_prefix}-backend-service"
  }
}

resource "aws_ecs_service" "frontend" {
  name            = "${local.name_prefix}-frontend"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.frontend.arn
  desired_count   = var.frontend_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.ecs_security_group_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.frontend.arn
    container_name   = "frontend"
    container_port   = var.container_port
  }

  depends_on = [aws_lb_listener.https]

  tags = {
    Name = "${local.name_prefix}-frontend-service"
  }
}

resource "aws_wafv2_web_acl" "alb" {
  #checkov:skip=CKV2_AWS_31:Full WAF request logging is disabled for PoC cost and PHI/header-minimization; ALB/CloudTrail/application audit logs remain.
  name  = "${local.name_prefix}-alb-waf"
  scope = "REGIONAL"

  default_action {
    allow {}
  }

  rule {
    name     = "aws-common-rules"
    priority = 1

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${local.name_prefix}-common-rules"
      sampled_requests_enabled   = false
    }
  }

  rule {
    name     = "aws-known-bad-inputs"
    priority = 2

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesKnownBadInputsRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${local.name_prefix}-known-bad-inputs"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "${local.name_prefix}-alb-waf"
    sampled_requests_enabled   = false
  }

  tags = {
    Name = "${local.name_prefix}-alb-waf"
  }
}

resource "aws_wafv2_web_acl_association" "alb" {
  resource_arn = aws_lb.main.arn
  web_acl_arn  = aws_wafv2_web_acl.alb.arn
}
