output "medical_images_bucket_name" {
  description = "S3 bucket name for medical images"
  value       = aws_s3_bucket.medical_images.id
}

output "medical_images_bucket_arn" {
  description = "S3 bucket ARN for medical images"
  value       = aws_s3_bucket.medical_images.arn
}

output "audit_logs_bucket_name" {
  description = "S3 bucket name for audit logs"
  value       = aws_s3_bucket.audit_logs.id
}

output "audit_logs_bucket_arn" {
  description = "S3 bucket ARN for audit logs"
  value       = aws_s3_bucket.audit_logs.arn
}
