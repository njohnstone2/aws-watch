resource "aws_cloudwatch_event_rule" "event_forwarding" {
  provider = aws.source_region

  name        = "${var.source_region}_event_forwarding"
  description = "Forwards Eventbridge events to a specified region"
  tags        = var.tags

  event_pattern = <<EOF
{
  "detail": {
    "eventSource": ["iam.amazonaws.com"],
    "eventName": ["CreatePolicy", "CreateRole", "CreateUser"]
  },
  "source": ["aws.iam"]
}
EOF
}

resource "aws_cloudwatch_event_target" "eventbridge_external" {
  provider = aws.source_region

  rule      = aws_cloudwatch_event_rule.event_forwarding.name
  target_id = "EventBridgeForwarding"
  arn       = "arn:aws:events:${var.target_region}:${var.aws_account_id}:event-bus/${var.event_bus_name}"
  role_arn  = aws_iam_role.eventbridge_remote_invoke.arn
}
