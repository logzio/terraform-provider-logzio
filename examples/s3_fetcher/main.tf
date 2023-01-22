variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_s3_fetcher" "my_s3_fetcher_keys" {
  aws_access_key = "my_access_key"
  aws_secret_key = "my_secret_key"
  bucket_name = "my_bucket"
  active = true
  add_s3_object_key_as_log_field = false
  aws_region = "EU_WEST_3"
  logs_type = "S3Access"
}

resource "logzio_s3_fetcher" "my_s3_fetcher_arn" {
  aws_arn = "my_arn"
  bucket_name = "my_bucket"
  active = true
  add_s3_object_key_as_log_field = false
  aws_region = "AP_SOUTHEAST_2"
  logs_type = "cloudfront"
}