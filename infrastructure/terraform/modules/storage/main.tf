locals {
  name_prefix   = "${var.project_name}-${var.environment}"
  bucket_prefix = "${var.project_name}-${var.environment}-${data.aws_caller_identity.current.account_id}"
}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "audit_logs" {
  statement {
    sid     = "AWSCloudTrailAclCheck"
    actions = ["s3:GetBucketAcl"]

    principals {
      type        = "Service"
      identifiers = ["cloudtrail.amazonaws.com"]
    }

    resources = [aws_s3_bucket.audit_logs.arn]
  }

  statement {
    sid     = "AWSCloudTrailWrite"
    actions = ["s3:PutObject"]

    principals {
      type        = "Service"
      identifiers = ["cloudtrail.amazonaws.com"]
    }

    resources = ["${aws_s3_bucket.audit_logs.arn}/AWSLogs/${data.aws_caller_identity.current.account_id}/*"]

    condition {
      test     = "StringEquals"
      variable = "s3:x-amz-acl"
      values   = ["bucket-owner-full-control"]
    }
  }

  statement {
    sid     = "AWSLoadBalancerWrite"
    actions = ["s3:PutObject"]

    principals {
      type        = "Service"
      identifiers = ["logdelivery.elasticloadbalancing.amazonaws.com"]
    }

    resources = ["${aws_s3_bucket.audit_logs.arn}/alb/AWSLogs/${data.aws_caller_identity.current.account_id}/*"]
  }

  statement {
    sid     = "S3ServerAccessLogsWrite"
    actions = ["s3:PutObject"]

    principals {
      type        = "Service"
      identifiers = ["logging.s3.amazonaws.com"]
    }

    resources = ["${aws_s3_bucket.audit_logs.arn}/s3-access/*"]
  }
}

resource "aws_s3_bucket" "medical_images" {
  #checkov:skip=CKV_AWS_144:Cross-region replication is documented as future DR work for this PoC.
  #checkov:skip=CKV2_AWS_62:S3 event notifications are not required for current medical image flow.
  bucket = "${local.bucket_prefix}-medical-images"

  tags = {
    Name        = "${local.name_prefix}-medical-images"
    DataClass   = "phi"
    Description = "Medical image object storage"
  }
}

resource "aws_s3_bucket_public_access_block" "medical_images" {
  bucket = aws_s3_bucket.medical_images.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "medical_images" {
  bucket = aws_s3_bucket.medical_images.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = var.kms_key_arn
      sse_algorithm     = "aws:kms"
    }
  }
}

resource "aws_s3_bucket_logging" "medical_images" {
  bucket = aws_s3_bucket.medical_images.id

  target_bucket = aws_s3_bucket.audit_logs.id
  target_prefix = "s3-access/medical-images/"
}

resource "aws_s3_bucket_versioning" "medical_images" {
  bucket = aws_s3_bucket.medical_images.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "medical_images" {
  bucket = aws_s3_bucket.medical_images.id

  rule {
    id     = "abort-incomplete-multipart-uploads"
    status = "Enabled"

    filter {}

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }
}

resource "aws_s3_bucket" "audit_logs" {
  #checkov:skip=CKV_AWS_18:This bucket is the centralized log destination; enabling access logs on itself creates recursive logs.
  #checkov:skip=CKV_AWS_144:Cross-region replication is documented as future DR work for this PoC.
  #checkov:skip=CKV2_AWS_62:S3 event notifications are not required for centralized audit log storage.
  bucket = "${local.bucket_prefix}-audit-logs"

  tags = {
    Name        = "${local.name_prefix}-audit-logs"
    DataClass   = "audit"
    Description = "AWS audit and compliance logs"
  }
}

resource "aws_s3_bucket_public_access_block" "audit_logs" {
  bucket = aws_s3_bucket.audit_logs.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "audit_logs" {
  bucket = aws_s3_bucket.audit_logs.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = var.kms_key_arn
      sse_algorithm     = "aws:kms"
    }
  }
}

resource "aws_s3_bucket_versioning" "audit_logs" {
  bucket = aws_s3_bucket.audit_logs.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "audit_logs" {
  bucket = aws_s3_bucket.audit_logs.id

  rule {
    id     = "retain-audit-logs"
    status = "Enabled"

    filter {}

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }

    expiration {
      days = 2555
    }
  }
}

resource "aws_s3_bucket_policy" "audit_logs" {
  bucket = aws_s3_bucket.audit_logs.id
  policy = data.aws_iam_policy_document.audit_logs.json
}
