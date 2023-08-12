resource "aws_ecr_pull_through_cache_rule" "ecr_public" {
  count                 = var.create_ecr_pullthrough_cache ? 1 : 0
  ecr_repository_prefix = "ecr-public"
  upstream_registry_url = "public.ecr.aws"
}
