resource "aws_lambda_function" "rds_creds_rotator" {
  function_name = "rds-rotator"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_binaries.key

  runtime = "python3.8"
  handler = "rds_rotator.lambda_handler"

  source_code_hash = data.archive_file.payload.output_base64sha256

  role = aws_iam_role.rds_rotator.arn

  layers = [
    aws_lambda_layer_version.pymysql.arn
  ]

  vpc_config {
    security_group_ids = [aws_security_group.rds.id]
    subnet_ids         = module.vpc.private_subnets
  }

  environment {
    variables = {
      "SECRETS_MANAGER_ENDPOINT" = "https://secretsmanager.${data.aws_region.current.name}.amazonaws.com"
    }
  }
}

resource "aws_lambda_permission" "secrets_manager_trigger" {
  function_name = aws_lambda_function.rds_creds_rotator.function_name

  statement_id_prefix = "allowSecretsManagerTrigger"

  action = "lambda:InvokeFunction"
  principal     = "secretsmanager.amazonaws.com"
  source_account = data.aws_caller_identity.current.account_id
}

resource "aws_lambda_layer_version" "pymysql" {
  filename   = "files/pymysql.zip"
  layer_name = "pymysql"

  compatible_runtimes = ["python3.8"]
}
