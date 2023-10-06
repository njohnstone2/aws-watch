resource "aws_secretsmanager_secret" "slack_token" {
  name                    = "slack_token"
  description             = "The slack bot token. Starts with xoxb"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "slack_token_value" {
  secret_id     = aws_secretsmanager_secret.slack_token.id
  secret_string = "<ENTER_TOKEN>"
}

resource "aws_secretsmanager_secret" "slack_channel_id" {
  name                    = "slack_channel_id"
  description             = "the ID of the slack channel messages will be sent to"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "slack_channel_id_value" {
  secret_id     = aws_secretsmanager_secret.slack_channel_id.id
  secret_string = "<ENTER_CHANNEL_ID>"
}
