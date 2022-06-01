resource "aws_iam_role" "rds_rotator" {
  name = "rds-rotator"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

#
# Basic Execution role attachment
resource "aws_iam_role_policy_attachment" "rds_rotator_basic" {
  role       = aws_iam_role.rds_rotator.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data aws_iam_policy_document "secrets_manager_rotation" {
  statement {
    effect = "Allow"
    actions = [
      "secretsmanager:DescribeSecret",
      "secretsmanager:GetSecretValue",
      "secretsmanager:PutSecretValue",
      "secretsmanager:UpdateSecretVersionStage"
    ]

    resources = [data.aws_secretsmanager_secret.rds.0.arn]
  }

  statement {
    effect = "Allow"
    actions = ["secretsmanager:GetRandomPassword"]
    resources = ["*"]
  }

  statement {
    effect = "Allow"
    actions = [
      "ec2:CreateNetworkInterface",
      "ec2:DeleteNetworkInterface",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DetachNetworkInterface"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "rds_secrets_manager_rotator" {
  name_prefix = "secrets-manager-rds-rotator-lambda"
  policy = data.aws_iam_policy_document.secrets_manager_rotation.json
}

resource "aws_iam_role_policy_attachment" "rds_rotator_secrets_manager" {
  role = aws_iam_role.rds_rotator.name
  policy_arn = aws_iam_policy.rds_secrets_manager_rotator.arn
}
