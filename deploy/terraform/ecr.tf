resource "aws_ecr_pull_through_cache_rule" "ecr_public" {
  count                 = var.create_ecr_pullthrough_cache ? 1 : 0
  ecr_repository_prefix = "ecr-public"
  upstream_registry_url = "public.ecr.aws"
}

resource "aws_ecr_repository" "aws_watch" {
  name                 = "aws-watch"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = {
    Name = "aws-watch"
  }
}

resource "aws_ecr_lifecycle_policy" "aws_watch" {
  repository = aws_ecr_repository.aws_watch.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep the most recent 20 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 20
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

output "aws_watch_ecr_repository_url" {
  description = "Private ECR repository URL for aws-watch"
  value       = aws_ecr_repository.aws_watch.repository_url
}
