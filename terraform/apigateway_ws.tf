data "aws_iam_policy_document" "backend" {
  statement {
    sid = "1"

    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]

    resources = [
      "arn:aws:s3:::*",
    ]
  }

  statement {
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      "arn:aws:s3:::${var.s3_bucket_name}",
    ]

    condition {
      test     = "StringLike"
      variable = "s3:prefix"

      values = [
        "",
        "home/",
        "home/&{aws:username}/",
      ]
    }
  }

  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "arn:aws:s3:::${var.s3_bucket_name}/home/&{aws:username}",
      "arn:aws:s3:::${var.s3_bucket_name}/home/&{aws:username}/*",
    ]
  }
}

resource "aws_iam_policy" "backend" {
  name   = "backend_policy"
  path   = "/"
  policy = data.aws_iam_policy_document.backend.json
}

resource "aws_iam_role" "backend_role" {
  name = "backend_role"

  assume_role_policy = jsondecode({
    Version = "2024-02-06"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_apigatewayv2_api" "backend" {
  name                       = "backend-websocket-api"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_integration" "backend" {
  api_id           = aws_apigatewayv2_api.backend
  integration_type = "MOCK"
}

resource "aws_apigatewayv2_integration_response" "backend" {
  api_id                   = aws_apigatewayv2_api.backend.id
  integration_id           = aws_apigatewayv2_integration.backend.id
  integration_response_key = "/200/"
}

resource "aws_apigatewayv2_route" "connect" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "$default"
}

resource "aws_apigatewayv2_route" "message" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "$default"
}

resource "aws_apigatewayv2_route" "disconnect" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "$default"
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
  api_id = aws_apigatewayv2_api.backend.id
  name   = "backend-stage"
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.backend.function_name
  principal     = "events.amazonaws.com"
}