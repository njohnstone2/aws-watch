# Example values only. These need to be substituted appropriately with your account
aws_account_id            = 123456789012
aws_region                = "us-east-1"
cloudwatch_log_group_name = "CloudTrailManagementEvents"

event_forwarding_enabled  = true
events_cloudtrail_enabled = true

s3_config_content  = <<EOF
subscribers:
  - id: teamA
    name: Team A
    notifiers:
      grafana-oncall:
        enabled: true
        webhook-url: http://localhost:8080/integrations/v1/formatted_webhook/ABCDEFGHIJKL1234567890XYZ/
        sources:
          - "iam"
      slack:
        enabled: true
        channel-id: AAAAAAAAAAA
        sources:
          - "*"
  - id: teamB
    name: Team B
    notifiers:
      grafana-oncall:
        enabled: false
      slack:
        enabled: true
        channel-id: BBBBBBBBBBB
        sources:
          - "*"

EOF
