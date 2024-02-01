/*
TODO: add the following resources

aws_dynamodb_table (x3)
*/

provider "aws" {
    region = "us-east-1"
}

resource "aws_dynamodb_table" "Group" {
  name           = "GroupTable"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "GroupID"

  attribute {
    name = "GroupID"
    type = "S"
  }
}

resource "aws_dynamodb_table" "User" {
  name           = "UserTable"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "UserID"

  attribute {
    name = "UserID"
    type = "S"
  }
}

resource "aws_dynamodb_table" "Message" {
  name           = "MessageTable"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "MessageId"
  range_key      = "Timestamp"

  attribute {
    name = "UserID"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "N"
  }
}

//ignore poll/pollresult