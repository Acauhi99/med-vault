output "vpc_flow_logs_group_name" {
  description = "VPC Flow Logs CloudWatch log group name"
  value       = aws_cloudwatch_log_group.vpc_flow_logs.name
}

output "cloudtrail_arn" {
  description = "CloudTrail ARN"
  value       = aws_cloudtrail.main.arn
}

output "config_logs_bucket_name" {
  description = "AWS Config delivery bucket name"
  value       = aws_s3_bucket.config_logs.id
}
