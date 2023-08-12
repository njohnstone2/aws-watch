module "events_cloudtrail" {
  count  = var.events_cloudtrail_enabled ? 1 : 0
  source = "./modules/events_cloudtrail"

  aws_account_id = data.aws_caller_identity.current.account_id
  aws_region     = var.aws_region
  retention_days = 1
  trail_name     = var.cloudwatch_log_group_name
}
