variable "aws_account_id" {
  type        = string
  description = "AWS Account ID"
}

variable "source_region" {
  type        = string
  description = "Source Eventbridge events to be forwarded to another region within the same account"
}

variable "target_region" {
  type        = string
  description = "Target region to receive Eventbridge events within the same account"
}

variable "event_bus_name" {
  type        = string
  description = "Name of the Eventbridge bus"
  default     = "default"
}

variable "tags" {
  type    = map(string)
  default = {}
}
