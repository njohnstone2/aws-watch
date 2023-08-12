variable "aws_account_id" {
  type        = string
  description = "AWS Account ID"
}

variable "aws_region" {
  type        = string
  description = "AWS Region"
}

variable "retention_days" {
  type        = number
  description = "Number of days to retain events"
  default     = 1
}

variable "trail_name" {
  type        = string
  description = "The name of the Cloudtrail Trail"
  default     = "ManagementEvents"
}
