resource "aws_s3_bucket" "frontend" {
  bucket = "lemmeknow-frontend"
}

resource "aws_s3_bucket_ownership_controls" "frontend" {
  bucket = aws_s3_bucket.frontend.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "frontend" {
  depends_on = [
    aws_s3_bucket_ownership_controls.frontend,
    aws_s3_bucket_public_access_block.frontend,
  ]

  bucket = aws_s3_bucket.frontend.id
  acl    = "public-read"
}

resource "aws_s3_bucket_cors_configuration" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    allowed_headers = ["*"]
  }
}

data "aws_iam_policy_document" "frontend" {
  statement {
    sid       = "public"
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.frontend.arn}/*"]

    principals {
      type        = "*"
      identifiers = ["*"]
    }
  }
}

resource "aws_s3_bucket_policy" "frontend" {
  bucket = aws_s3_bucket.frontend.id
  policy = data.aws_iam_policy_document.frontend.json
}

resource "aws_s3_object" "frontend" {
  bucket        = aws_s3_bucket.frontend.bucket
  cache_control = "no-cache"
  content_type = {
    "html" = "text/html",
    "xml"  = "application/xml",
    "txt"  = "text/plain",
    "js"   = "application/javascript",
    "css"  = "text/css",
    "json" = "application/json",
    "ico"  = "image/x-icon",
    "png"  = "image/png",
  }[split(".", each.value)[length(split(".", each.value)) - 1]]
  etag     = filemd5("../frontend/build/${each.value}")
  for_each = fileset("../frontend/build", "**/*.*")
  key      = each.value
  source   = "../frontend/build/${each.value}"
}
