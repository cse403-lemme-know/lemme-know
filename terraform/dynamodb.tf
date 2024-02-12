resource "aws_dynamodb_table" "group" {
  name         = "lemmeknow-groups"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "GroupID"

  attribute {
    name = "GroupID"
    type = "N"
  }
}

resource "aws_dynamodb_table" "user" {
  name         = "lemmeknow-users"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserID"

  attribute {
    name = "UserID"
    type = "N"
  }
}

resource "aws_dynamodb_table" "message" {
  name         = "lemmeknow-messages"
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

resource "aws_dynamodb_table" "variable" {
  name         = "lemmeknow-variables"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Name"

  attribute {
    name = "Name"
    type = "S"
  }
}