resource "aws_api_gateway_rest_api" "backend" {
  binary_media_types = ["*/*"]
  endpoint_configuration {
    types = ["REGIONAL"]
  }
  description              = "LemmeKnow backend REST API"
  minimum_compression_size = 1000
  name                     = "lemmeknow-backend-rest"
}

resource "aws_api_gateway_resource" "backend_rest_proxy" {
  parent_id   = aws_api_gateway_rest_api.backend.root_resource_id
  path_part   = "{proxy+}"
  rest_api_id = aws_api_gateway_rest_api.backend.id
}

resource "aws_api_gateway_method" "backend_rest_root" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_rest_api.backend.root_resource_id
  rest_api_id   = aws_api_gateway_rest_api.backend.id
}

resource "aws_api_gateway_integration" "backend_rest_root" {
  http_method             = aws_api_gateway_method.backend_rest_root.http_method
  integration_http_method = "POST"
  rest_api_id             = aws_api_gateway_rest_api.backend.id
  resource_id             = aws_api_gateway_rest_api.backend.root_resource_id
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.backend.invoke_arn
}

resource "aws_api_gateway_method" "backend_rest_proxy" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_resource.backend_rest_proxy.id
  rest_api_id   = aws_api_gateway_rest_api.backend.id
}

resource "aws_api_gateway_integration" "backend_rest_proxy" {
  http_method             = aws_api_gateway_method.backend_rest_proxy.http_method
  integration_http_method = "POST"
  rest_api_id             = aws_api_gateway_rest_api.backend.id
  resource_id             = aws_api_gateway_method.backend_rest_proxy.resource_id
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.backend.invoke_arn
}

resource "aws_api_gateway_deployment" "backend_rest" {
  depends_on  = [aws_api_gateway_method.backend_rest_root, aws_api_gateway_method.backend_rest_proxy, aws_api_gateway_integration.backend_rest_root, aws_api_gateway_integration.backend_rest_proxy]
  rest_api_id = aws_api_gateway_rest_api.backend.id
  stage_name  = "prod"
}

resource "aws_lambda_permission" "backend_rest" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.backend.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.backend.execution_arn}/*/*"
  statement_id  = "APIGatewayRest"
}
