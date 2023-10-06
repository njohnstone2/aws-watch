resource "aws_cloudtrail" "management" {
  name           = var.trail_name
  s3_bucket_name = aws_s3_bucket.ct_management_logs.id

  include_global_service_events = true
  is_multi_region_trail         = true

  cloud_watch_logs_group_arn = "${aws_cloudwatch_log_group.ct_management_logs.arn}:*"
  cloud_watch_logs_role_arn  = aws_iam_role.ct_management_logs.arn

  advanced_event_selector {
    name = "Management events selector"

    field_selector {
      field  = "eventCategory"
      equals = ["Management"]
    }
  }

  depends_on = [aws_s3_bucket_policy.ct_management_logs]
}

resource "aws_cloudwatch_log_group" "ct_management_logs" {
  name = var.trail_name

  retention_in_days = var.retention_days
}
