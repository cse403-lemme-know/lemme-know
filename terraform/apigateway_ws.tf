data "aws_iam_policy_document" "backend_ws" {
  statement {
    sid = "gateway"
    actions = [
      "lambda:InvokeFunction",
    ]
    resources = [aws_lambda_function.backend.arn]
  }
}

resource "aws_iam_policy" "backend_ws" {
  name   = "lemmeknow-backend-ws"
  path   = "/"
  policy = data.aws_iam_policy_document.backend_ws.json
}

resource "aws_iam_role" "backend_ws" {
  name                = "lemmeknow-backend-ws"
  assume_role_policy  = data.aws_iam_policy_document.backend_ws_assume_role.json
  managed_policy_arns = [aws_iam_policy.backend_ws.arn]
}

data "aws_iam_policy_document" "backend_ws_assume_role" {
  statement {
    sid = "gateway"
    actions = [
      "sts:AssumeRole",
    ]
    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }
  }
}

resource "aws_apigatewayv2_api" "backend" {
  name                       = "lemmeknow-backend-ws"
  description                = "LemmeKnow backend WS API"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_integration" "backend_ws_proxy" {
  api_id                    = aws_apigatewayv2_api.backend.id
  content_handling_strategy = "CONVERT_TO_TEXT"
  credentials_arn           = aws_iam_role.backend_ws.arn
  integration_method        = "POST"
  integration_type          = "AWS_PROXY"
  integration_uri           = aws_lambda_function.backend.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_integration_response" "backend" {
  api_id                   = aws_apigatewayv2_api.backend.id
  integration_id           = aws_apigatewayv2_integration.backend_ws_proxy.id
  integration_response_key = "/200/"
}

resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.backend_ws_proxy.id}"
}

resource "aws_apigatewayv2_route" "connect" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "$connect"
  target    = "integrations/${aws_apigatewayv2_integration.backend_ws_proxy.id}"
}

resource "aws_apigatewayv2_route" "message" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "MESSAGE"
  target    = "integrations/${aws_apigatewayv2_integration.backend_ws_proxy.id}"
}

resource "aws_apigatewayv2_route" "disconnect" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "$disconnect"
  target    = "integrations/${aws_apigatewayv2_integration.backend_ws_proxy.id}"
}

resource "aws_apigatewayv2_route_response" "connect" {
  api_id             = aws_apigatewayv2_api.backend.id
  route_id           = aws_apigatewayv2_route.connect.id
  route_response_key = "$default"
}

resource "aws_apigatewayv2_route_response" "message" {
  api_id             = aws_apigatewayv2_api.backend.id
  route_id           = aws_apigatewayv2_route.message.id
  route_response_key = "$default"
}

resource "aws_apigatewayv2_route_response" "disconnect" {
  api_id             = aws_apigatewayv2_api.backend.id
  route_id           = aws_apigatewayv2_route.disconnect.id
  route_response_key = "$default"
}

resource "aws_apigatewayv2_stage" "backend" {
  api_id      = aws_apigatewayv2_api.backend.id
  name        = "backend-stage"
  auto_deploy = true
  default_route_settings {
    throttling_burst_limit = 64
    throttling_rate_limit  = 64
  }

}

resource "aws_lambda_permission" "backend_ws" {
  statement_id  = "APIGatewayWs"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.backend.function_name
  principal     = "events.amazonaws.com"
}