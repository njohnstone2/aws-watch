resource "random_string" "uniqueness" {
  length  = 16
  lower   = true
  upper   = true
  special = false
}

resource "aws_s3_bucket" "main" {
  bucket = "aws-watch-${lower(random_string.uniqueness.result)}"

  tags = {
    Name = "aws-watch"
  }
}

## TODO move to tfvars
resource "aws_s3_object" "config" {
  bucket  = aws_s3_bucket.main.id
  key     = var.s3_config_filename
  content = var.s3_config_content
}