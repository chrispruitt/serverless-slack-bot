resource "aws_apigatewayv2_api" "main" {
  name          = "${terraform.workspace}-${var.name}"
  description   = "This is an API to for the ${terraform.workspace} ${var.name}"
  protocol_type = "HTTP"
  depends_on    =   [aws_lambda_function.main]
}

resource "aws_apigatewayv2_deployment" "main" {
  depends_on = [
    aws_apigatewayv2_route.any_root,

  ]

  api_id = aws_apigatewayv2_api.main.id
  triggers = {
    "SourceCodeHash" = data.aws_s3_bucket_object.source_code_hash.body
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id                 = aws_apigatewayv2_api.main.id
  description            = "Lambda example"
  integration_type       = "AWS_PROXY"
  connection_type        = "INTERNET"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.main.invoke_arn
  passthrough_behavior   = "WHEN_NO_MATCH"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_stage" "all" {
  api_id = aws_apigatewayv2_api.main.id
  name   = "default"
  deployment_id = aws_apigatewayv2_deployment.main.id
}

resource "aws_lambda_permission" "apigw_get_config" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.main.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.main.execution_arn}/*/*"
}

// API Paths
resource "aws_apigatewayv2_route" "any_root" {
  api_id         = aws_apigatewayv2_api.main.id
  target         = "integrations/${aws_apigatewayv2_integration.lambda.id}"
  route_key      = "ANY /{proxy+}"
  operation_name = "proxy"
}

output "api_endpoint" {
  value = aws_apigatewayv2_api.main.api_endpoint
}
