terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  backend "http" {
    address        = "https://api.tfstate.dev/github/v1"
    lock_address   = "https://api.tfstate.dev/github/v1/lock"
    unlock_address = "https://api.tfstate.dev/github/v1/lock"
    lock_method    = "PUT"
    unlock_method  = "DELETE"
    username       = "cse403-lemmeknow/lemmeknow"
  }
}

provider "aws" {
  profile                  = "cse403"
  region                   = "us-east-1"
  shared_credentials_files = ["$HOME/.aws/credentials"]
}

data "aws_caller_identity" "current" {}