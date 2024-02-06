locals {
  frontend_origin_id     = "S3-frontend"
  backend_rest_origin_id = "API-GATEWAY-backend-rest"
  backend_ws_origin_id   = "API-GATEWAY-backend-ws"
}

resource "aws_cloudfront_origin_access_identity" "backend" {
  comment = "origin access identity"
}

resource "aws_cloudfront_distribution" "backend" {
  depends_on          = [aws_acm_certificate_validation.backend, aws_s3_bucket.frontend]
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  aliases             = [var.domain]

  origin {
    domain_name = aws_s3_bucket_website_configuration.frontend.website_endpoint
    origin_id   = local.frontend_origin_id
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.1", "TLSv1.2"]
    }
  }
  origin {
    domain_name = element(split("/", aws_api_gateway_deployment.backend.invoke_url), 2)
    origin_id   = local.backend_rest_origin_id
    origin_path = "/prod"
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }
  /*
  origin {
    domain_name = element(split("/", aws_apigatewayv2_stage.backend.invoke_url), 2)
    origin_id = local.backend_ws_origin_id
    origin_path = ""
    custom_origin_config {
      http_port = 80
      https_port = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols = ["TLSv1.2"]
    }
  }
  */
  default_cache_behavior {
    allowed_methods = ["GET", "HEAD", "OPTIONS"]
    cached_methods  = ["GET", "HEAD"]
    compress        = true
    forwarded_values {
      cookies {
        forward = "all"
      }
      query_string = false
    }
    target_origin_id       = local.frontend_origin_id
    viewer_protocol_policy = "redirect-to-https"
  }
  ordered_cache_behavior {
    allowed_methods = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    cached_methods  = ["GET", "HEAD"]
    compress        = true
    forwarded_values {
      cookies {
        forward = "all"
      }
      headers      = []
      query_string = true
    }
    max_ttl                = 600
    path_pattern           = "api/*"
    target_origin_id       = local.backend_rest_origin_id
    viewer_protocol_policy = "redirect-to-https"
  }
  ordered_cache_behavior {
    allowed_methods          = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    cache_policy_id          = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad"
    cached_methods           = ["GET", "HEAD"]
    compress                 = true
    origin_request_policy_id = "b689b0a8-53d0-40ab-baf2-68738e2966ac"
    path_pattern             = "ws/*"
    target_origin_id         = local.backend_ws_origin_id
    viewer_protocol_policy   = "redirect-to-https"
  }

  price_class = "PriceClass_100"

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

}

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

resource "aws_s3_bucket_policy" "frontend" {
  bucket     = aws_s3_bucket.frontend.id
  depends_on = [aws_cloudfront_origin_access_identity.main]
  policy     = data.aws_iam_policy_document.frontend.json
}

data "aws_iam_policy_document" "frontend" {
  statement {
    actions = ["s3:GetObject"]
    principals {
      type        = "*"
      identifiers = ["*"]
    }
    resources = ["${aws_s3_bucket.frontend.arn}/*"]
    sid       = "cloudfront"
  }
}
