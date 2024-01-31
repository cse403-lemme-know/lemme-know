resource "aws_cloudwatch_event_rule" "backend_cron" {
    name = "lemmeknow-backend-cron"
    description = "Run lemmeknow-backend-api lambda every hour"
    schedule_expression = "cron(0 * * * ? *)"
}

resource "aws_cloudwatch_event_target" "backend_target" {
    arn = aws_lambda_function.backend.arn
    rule = aws_cloudwatch_event_rule.backend_cron.name
    target_id = "lemmeknow-backend"
}

resource "aws_lambda_permission" "backend_cron" {
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.backend.function_name
    principal = "events.amazonaws.com"
    statement_id = "EventBridge"
}
