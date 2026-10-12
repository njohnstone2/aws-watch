terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"

      configuration_aliases = [
        aws.source_region,
        aws.target_region
      ]
    }
  }
}