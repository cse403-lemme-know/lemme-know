resource "aws_iam_user" "backend_dev" {
  name = "lemmeknow-backend-dev"
  path = "/system/"
}

data "aws_iam_policy_document" "backend_dev" {
  statement {
    effect    = "Allow"
    actions   = ["lambda:UpdateFunctionCode", "lambda:InvokeFunction"]
    resources = ["${aws_lambda_function.backend.arn}"]
  }
}

resource "aws_iam_user_policy" "backend_dev" {
  name   = "lemmeknow-backend-dev"
  user   = aws_iam_user.backend_dev.name
  policy = data.aws_iam_policy_document.backend_dev.json
}

resource "aws_iam_access_key" "backend_dev" {
  user = aws_iam_user.backend_dev.name
}

resource "local_file" "backend_dev_access_key_id" {
  content  = aws_iam_access_key.backend_dev.id
  filename = "../backend/access_key_id"
}

resource "local_file" "backend_dev_secret_access_key" {
  content  = aws_iam_access_key.backend_dev.secret
  filename = "../backend/secret_access_key"
}