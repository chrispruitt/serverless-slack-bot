data "aws_s3_bucket_object" "source_code_hash" {
  bucket = var.s3_bucket
  key    = var.s3_key
}

resource "aws_cloudwatch_log_group" "main" {
  name              = "/aws/lambda/${aws_lambda_function.main.function_name}"
  retention_in_days = 30
}

resource "aws_lambda_function" "main" {
  s3_bucket        = var.s3_bucket
  s3_key           = var.s3_key
  role             = aws_iam_role.main.arn
  function_name    = "${terraform.workspace}-${var.name}"
  handler          = "serverless-slack-bot"
  runtime          = "go1.x"
  publish          = true
  source_code_hash = data.aws_s3_bucket_object.source_code_hash.body

  environment {
    variables = {
      SLACK_OAUTH_ACCESS_TOKEN = var.slack_oauth_access_token
      SLACK_VERIFICATION_TOKEN = var.slack_verification_token
    }
  }

  lifecycle {
    create_before_destroy = false
  }
}
