resource "logzio_archive_logs" "test_to_datasource" {
  storage_type = "S3"
  compressed = false
  amazon_s3_storage_settings {
    credentials_type = "IAM"
    s3_path = "%s"
    s3_iam_credentials_arn = "%s"
  }
}

data "logzio_archive_logs" "my_archive_datasource" {
  archive_id = "${logzio_archive_logs.test_to_datasource.id}"
  depends_on = ["logzio_archive_logs.test_to_datasource"]
}
