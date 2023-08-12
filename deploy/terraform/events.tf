module "event_forwarding" {
  count  = var.event_forwarding_enabled ? 1 : 0
  source = "./modules/event_forwarding"

  aws_account_id = data.aws_caller_identity.current.account_id
  source_region  = var.event_forwarding_source_region
  target_region  = var.event_forwarding_target_region

  providers = {
    aws.source_region = aws.source_region
    aws.target_region = aws
  }
}

provider "aws" {
  region = "us-east-1"
  alias  = "source_region"

  default_tags {
    tags = {
      Environment = "Management"
      ManagedBy   = "Terraform"
    }
  }
}
