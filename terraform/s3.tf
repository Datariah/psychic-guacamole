resource "aws_s3_bucket" "lambda_bucket" {
  bucket_prefix = "psychic-guacamole-lambdas-"

  force_destroy = true
}

data "archive_file" "payload" {
  type = "zip"

  source_file = "files/rds_rotator.py"
  output_path = "${path.module}/payload.zip"

}

resource "aws_s3_object" "lambda_binaries" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "payload.zip"
  source = data.archive_file.payload.output_path

  source_hash = filemd5(data.archive_file.payload.output_path)
}
