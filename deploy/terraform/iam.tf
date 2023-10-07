data "aws_caller_identity" "current" {}

resource "aws_iam_role" "execution_role" {
  name = var.app_name
  path = "/service-role/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cloudwatch_logs" {
  role       = aws_iam_role.execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "read_secrets" {
  role       = aws_iam_role.execution_role.name
  policy_arn = aws_iam_policy.read_secrets.arn
}

resource "aws_iam_policy" "read_secrets" {
  name        = "${var.app_name}-read-secrets"
  description = "Grants access to aws-watch secrets"

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ReadSecrets",
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret"
            ],
            "Resource": "arn:aws:secretsmanager:${var.aws_region}:${var.aws_account_id}:secret:*"
        }
    ]
}
POLICY
}
