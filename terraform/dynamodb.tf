resource "aws_dynamodb_table" "group" {
  name         = "GroupTable"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "GroupID"

  attribute {
    name = "GroupID"
    type = "N"
  }
}

resource "aws_dynamodb_table" "user" {
  name         = "UserTable"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserID"

  attribute {
    name = "UserID"
    type = "N"
  }
}

resource "aws_dynamodb_table" "message" {
  name         = "MessageTable"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "MessageId"
  range_key    = "Timestamp"

  attribute {
    name = "MessageId"
    type = "N"
  }

  attribute {
    name = "Timestamp"
    type = "N"
  }
}