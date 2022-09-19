resource "logzio_archive_logs" "my_s3_archive" {
  storage_type = "S3"
  enabled = false
  credentials_type = "IAM"
  s3_path = "some-path"
  s3_iam_credentials_arn = "some-arn"
}

resource "logzio_archive_logs" "my_s3_archive_keys" {
  storage_type = "S3"
  enabled = false
  credentials_type = "KEYS"
  s3_path = "some-path"
  aws_access_key = "some-access-key"
  aws_secret_key = "some-secret-key"
}

resource "logzio_archive_logs" "my_blob_archive" {
  storage_type = "BLOB"
  enabled = false
  tenant_id = "some-tenant-id"
  client_id = "some-client-id"
  client_secret = "some-client-secret"
  account_name = "some-account-name"
  container_name = "some-container-name"
}