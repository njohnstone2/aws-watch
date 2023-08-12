variable "aws_account_id" {
  type        = string
  description = "AWS Account ID"
}

variable "aws_region" {
  type        = string
  description = "AWS Region"
}

variable "app_name" {
  type        = string
  description = "Name of the application"
  default     = "aws-watch"
}

variable "cloudwatch_log_group_name" {
  type        = string
  description = "Cloudwatch log group containing cloudtrail events"
}

variable "image_repo" {
  type        = string
  description = "Docker image repository name used on the lambda function (e.g. example/aws-watch)"
  default     = "b3g7u0g4"
}

variable "image_tag" {
  type        = string
  description = "Docker image tag used on the lambda function"
  default     = "latest"
}

variable "create_ecr_pullthrough_cache" {
  type        = bool
  description = "Whether an ECR pull through cache should be created to track upstream project changes"
  default     = true
}

# Quickstart Modules
## Cloudtrail Management events
variable "enable_monitor_management_events" {
  type        = bool
  description = "Whether the Cloudtrail Management events module should be created"
  default     = false
}

## Region Event Forwarding
variable "event_forwarding_enabled" {
  type        = bool
  description = "Whether the Region event forwarding module should be created"
  default     = false
}

variable "event_forwarding_source_region" {
  type        = string
  description = "The source region for the event forwarding module"
  default     = "us-east-1"
}

variable "event_forwarding_target_region" {
  type        = string
  description = "The target region for the event forwarding module"
  default     = "ap-southeast-2"
}

variable "events_cloudtrail_enabled" {
  type        = bool
  description = "Whether the cloudtrail management event resources should be created"
  default     = false
}
