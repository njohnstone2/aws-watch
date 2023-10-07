resource "aws_lambda_function" "main" {
  function_name = var.app_name
  role          = aws_iam_role.execution_role.arn
  timeout       = "360"
  memory_size   = "128"

  package_type = "Image"
  image_uri    = var.create_ecr_pullthrough_cache ? "${var.aws_account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/ecr-public/${var.image_repo}/${var.app_name}:${var.image_tag}" : "${var.image_repo}:${var.image_tag}"

  environment {
    variables = {
      LOG_LEVEL = "info"
      REGION = var.aws_region
    }
  }
}

resource "aws_lambda_permission" "cloudtrail_invoke" {
  statement_id  = "AllowExecutionFromCloudtrail"
  action        = "lambda:InvokeFunction"
  function_name = var.app_name
  principal     = "logs.${var.aws_region}.amazonaws.com"
  source_arn    = "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:log-group:${var.cloudwatch_log_group_name}:*"
}

resource "aws_cloudwatch_log_subscription_filter" "cloudtrail_trigger" {
  name            = "Cloudtrail_CloudwatchLogTriggerIAM"
  log_group_name  = var.cloudwatch_log_group_name
  filter_pattern  = "{ $.eventSource != \"logs.amazonaws.com\" && ($.eventName = \"Create*\" || $.eventName = \"Update*\" || $.eventName = \"Delete*\") }"
  destination_arn = aws_lambda_function.main.arn
}
