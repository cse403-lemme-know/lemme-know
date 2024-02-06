resource "aws_route53_zone" "backend" {
  name = var.domain
}

// TODO: Uncomment when Cloudfront is added
/*
resource "aws_route53_record" "cloudfront" {
  alias {
    evaluate_target_health = false
    name                   = aws_cloudfront_distribution.backend.domain_name
    zone_id                = aws_cloudfront_distribution.backend.hosted_zone_id
  }
  name    = aws_route53_zone.backend.name
  type    = "A"
  zone_id = aws_route53_zone.backend.zone_id
}
*/