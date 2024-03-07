resource "aws_lambda_function" "backend" {
  filename         = "../backend/bin/bootstrap.zip"
  function_name    = "lemmeknow-backend"
  handler          = "bootstrap"
  memory_size      = 128
  role             = aws_iam_role.backend_role.arn
  runtime          = "provided.al2023"
  source_code_hash = filebase64sha256("../backend/bin/bootstrap.zip")
  timeout          = 10
  environment {
    variables = {
      AWS_API_GATEWAY_WS_ENDPOINT = replace(aws_apigatewayv2_stage.backend.invoke_url, "wss:", "https:")
      AWS_LAMBDA_ARN              = "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:lemmeknow-backend"
      AWS_SCHEDULER_ROLE_ARN      = aws_iam_role.backend_role.arn
      #AWS_REGION                 = var.region
    }
  }

}

resource "aws_iam_role" "backend_role" {
  name               = "lemmeknow-backend-role"
  assume_role_policy = data.aws_iam_policy_document.backend_role.json
}

data "aws_iam_policy_document" "backend_role" {
  statement {
    sid = "lambda"
    actions = [
      "sts:AssumeRole"
    ]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com", "scheduler.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "backend_policy" {
  name   = "lemmeknow-backend-policy"
  role   = aws_iam_role.backend_role.id
  policy = data.aws_iam_policy_document.backend_policy.json
}

data "aws_iam_policy_document" "backend_policy" {
  statement {
    sid = "logs"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:*:*:*"
    ]
  }
  statement {
    sid = "dynamodb"
    actions = [
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:UpdateItem"
    ]
    resources = [
      aws_dynamodb_table.user.arn,
      aws_dynamodb_table.group.arn,
      aws_dynamodb_table.message.arn,
      aws_dynamodb_table.connection.arn,
      aws_dynamodb_table.variable.arn
    ]
  }
  statement {
    sid = "websocket"
    actions = [
      "execute-api:*"
    ]
    resources = [
      "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.backend.id}/*/*/*"
    ]
  }
  statement {
    sid = "scheduler"
    actions = [
      "scheduler:CreateSchedule"
    ]
    resources = ["*"]
  }
  statement {
    sid = "iam"
    actions = [
      "iam:PassRole"
    ]
    resources = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/lemmeknow-backend-role"]
  }
}