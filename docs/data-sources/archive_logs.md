# Archive Logs Datasource

Provides a Logz.io archive logs resource. This can be used to create and manage Logz.io archive logs settings.

* Learn more about archive logs in the [Logz.io Docs](https://docs.logz.io/api/#tag/Archive-logs)

## Argument Reference

* `archive_id` - (String) Archive ID in the Logz.io database.

## Attribute Reference

* `storage_type` - (String) Specifies the storage provider.
* `enabled` - (Boolean) If `true`, archiving is currently enabled.
* `compressed` - (Boolean) If `true`, logs are compressed before they are archived.
* `aws_credentials_type` - (String) Applicable when `storage_type` is `S3`. Specifies which credentials will be used for authentication.
* `aws_s3_path` - (String) Applicable when `storage_type` is `S3`. Specify a path to the **root** of the S3 bucket.
* `aws_s3_iam_credentials_arn` - (String) Applicable when `storage_type` is `S3`. Amazon Resource Name (ARN) to uniquely identify the S3 bucket.
* `aws_access_key` - (String) Applicable when `storage_type` is `S3`. AWS access key.
* `azure_tenant_id` - (String) Applicable when `storage_type` is `BLOB`. Azure Directory (tenant) ID. The Tenant ID of the AD app.
* `azure_client_id` - (String) Applicable when `storage_type` is `BLOB`. Azure application (client) ID. The Client ID of the AD app.
* `azure_account_name` - (String) Applicable when `storage_type` is `BLOB`. Azure Storage account name.
* `azure_container_name` - (String) Applicable when `storage_type` is `BLOB`. Name of the container in the Storage account.
* `azure_blob_path` - (String) Optional virtual sub-folder specifying a path within the container.


