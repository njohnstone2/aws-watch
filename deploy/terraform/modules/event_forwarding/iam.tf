resource "aws_iam_role" "eventbridge_remote_invoke" {
  name = "${var.source_region}_event_forwarding"
  tags = var.tags

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "events.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "event_forwarding" {
  role       = aws_iam_role.eventbridge_remote_invoke.name
  policy_arn = aws_iam_policy.event_forwarding.arn
}

resource "aws_iam_policy" "event_forwarding" {
  name        = aws_iam_role.eventbridge_remote_invoke.name
  description = "Grants access to forward eventbridge events between AWS regions"
  tags        = var.tags

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "events:PutEvents",
      "Resource": "arn:aws:events:${var.target_region}:${var.aws_account_id}:event-bus/${var.event_bus_name}"
    }
  ]
}
POLICY
}
