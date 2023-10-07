resource "aws_s3_bucket" "ct_management_logs" {
  bucket        = "cloudtrail-management-events-${var.aws_account_id}"
  force_destroy = true

  tags = {
    Name = "cloudtrail-management-events-${var.aws_account_id}"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "ct_expire_objects" {
  bucket = aws_s3_bucket.ct_management_logs.id

  rule {
    id     = "expire-objects"
    status = "Enabled"

    expiration {
      days = var.retention_days
    }
  }
}

resource "aws_s3_bucket_policy" "ct_management_logs" {
  bucket = aws_s3_bucket.ct_management_logs.id
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AWSCloudTrailAclCheck20150319",
            "Effect": "Allow",
            "Principal": {
                "Service": "cloudtrail.amazonaws.com"
            },
            "Action": "s3:GetBucketAcl",
            "Resource": "arn:aws:s3:::${aws_s3_bucket.ct_management_logs.id}",
            "Condition": {
                "StringEquals": {
                    "AWS:SourceArn": "arn:aws:cloudtrail:${var.aws_region}:${var.aws_account_id}:trail/${var.trail_name}"
                }
            }
        },
        {
            "Sid": "AWSCloudTrailWrite20150319",
            "Effect": "Allow",
            "Principal": {
                "Service": "cloudtrail.amazonaws.com"
            },
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::${aws_s3_bucket.ct_management_logs.id}/AWSLogs/${var.aws_account_id}/*",
            "Condition": {
                "StringEquals": {
                    "s3:x-amz-acl": "bucket-owner-full-control",
                    "AWS:SourceArn": "arn:aws:cloudtrail:${var.aws_region}:${var.aws_account_id}:trail/${var.trail_name}"
                }
            }
        }
    ]
}
EOF
}
