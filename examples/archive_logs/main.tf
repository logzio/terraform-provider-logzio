resource "logzio_archive_logs" "my_s3_archive" {
  storage_type = "S3"
  enabled = false
  amazon_s3_storage_settings {
    credentials_type = "IAM"
    s3_path = "some-path"
    s3_iam_credentials_arn = "some-arn"
  }
}

resource "logzio_archive_logs" "my_blob_archive" {
  storage_type = "BLOB"
  azure_blob_storage_settings {
    tenant_id = "some-tenant-id"
    client_id = "some-client-id"
    client_secret = "some-client-secret"
    account_name = "some-account-name"
    container_name = "some-container-name"
  }
}