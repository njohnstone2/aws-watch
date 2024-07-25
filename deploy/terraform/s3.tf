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
  key     = "config.yaml"
  content = <<EOF
subscribers:
  - id: teamA
    name: Team A
    notifiers:
      grafana-oncall:
        enabled: false
        webhook-url: http://localhost:8080
        sources:
          - ec2
      slack:
        enabled: true
        channel-id: C056KLMHBL5
        sources:
          - "iam"
  - id: teamB
    name: Team B
    notifiers:
      grafana-oncall:
        enabled: false
      slack:
        enabled: true
        channel-id: C03SR20D8QL
        sources:
          - "*"

EOF
}