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
    domain_name = element(split("/", aws_api_gateway_deployment.backend_rest.invoke_url), 2)
    origin_id   = local.backend_rest_origin_id
    origin_path = "/prod"
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }
  origin {
    domain_name = element(split("/", aws_apigatewayv2_stage.backend.invoke_url), 2)
    origin_id   = local.backend_ws_origin_id
    origin_path = ""
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }
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

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate.backend.arn
    minimum_protocol_version = "TLSv1.1_2016"
    ssl_support_method       = "sni-only"
  }
}

