resource "aws_secretsmanager_secret" "slack_token" {
  name                    = "slack_token"
  description             = "The slack bot token. Starts with xoxb"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "slack_token_value" {
  secret_id     = aws_secretsmanager_secret.slack_token.id
  secret_string = "<ENTER_TOKEN>"
}
