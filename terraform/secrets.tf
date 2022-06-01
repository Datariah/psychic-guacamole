data "aws_secretsmanager_secret" "rds" {
  count = 0
  name = "psychic-guacamole-rds"
}

resource "aws_secretsmanager_secret_rotation" "rds_secret" {
  count = 0 # change me to 1 once db has been created and you've created your user
  secret_id           = data.aws_secretsmanager_secret.rds.0.id
  rotation_lambda_arn = aws_lambda_function.rds_creds_rotator.arn

  rotation_rules {
    automatically_after_days = 1
  }
}
