resource "logzio_archive_logs" "my_s3_archive" {
  storage_type = "S3"
  enabled = false
  aws_credentials_type = "IAM"
  aws_s3_path = "some-path"
  aws_s3_iam_credentials_arn = "some-arn"
}

resource "logzio_archive_logs" "my_s3_archive_keys" {
  storage_type = "S3"
  enabled = false
  aws_credentials_type = "KEYS"
  aws_s3_path = "some-path"
  aws_access_key = "some-access-key"
  aws_secret_key = "some-secret-key"
}

resource "logzio_archive_logs" "my_blob_archive" {
  storage_type = "BLOB"
  enabled = false
  azure_tenant_id = "some-tenant-id"
  azure_client_id = "some-client-id"
  azure_client_secret = "some-client-secret"
  azure_account_name = "some-account-name"
  azure_container_name = "some-container-name"
}