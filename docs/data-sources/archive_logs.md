# Archive Logs Resource

Provides a Logz.io archive logs resource. This can be used to create and manage Logz.io archive logs settings.

* Learn more about archive logs in the [Logz.io Docs](https://docs.logz.io/api/#tag/Archive-logs)

## Argument Reference
* `archive_id` - (String) Archive ID in the Logz.io database.

## Attribute Reference
* `storage_type` - (String) Specifies the storage provider.
  If `S3`, the `amazon_s3_storage_settings` are relevant.
  If `BLOB`, the `azure_blob_storage_settings` are relevant.
* `amazon_s3_storage_settings` - (Object) Applicable settings when the `storage_type` is `S3`.
* `azure_blob_storage_settings` - (Object) Applicable settings when the `storage_type` is `BLOB`.
* `enabled` - (Boolean) If `true`, archiving is currently enabled.
* `compressed` - (Boolean) If `true`, logs are compressed before they are archived.
* `amazon_s3_storage_settings.credentials_type` - (String) Specifies which credentials will be used for authentication.
* `amazon_s3_storage_settings.path` - (String) Specify a path to the **root** of the S3 bucket.
* `amazon_s3_storage_settings.s3_secret_credentials` - (Object) Applicable settings when the `credentials_type` is `KEYS`.
* `amazon_s3_storage_settings.s3_iam_credentials_arn` - (String) Amazon Resource Name (ARN) to uniquely identify the S3 bucket.
* `amazon_s3_storage_settings.s3_external_id` - (String) The external id that gives Logz.io access to your S3 bucket.
* `amazon_s3_storage_settings.s3_secret_credentials.access_key` - (String) AWS access key.
* `amazon_s3_storage_settings.s3_secret_credentials.secret_key` - (String) AWS secret key.
* `azure_blob_storage_settings.tenant_id` - (String) Azure Directory (tenant) ID. The Tenant ID of the AD app.
* `azure_blob_storage_settings.client_id` - (String) Azure application (client) ID. The Client ID of the AD app.
* `azure_blob_storage_settings.client_secret` - (String) Azure client secret.
* `azure_blob_storage_settings.account_name` - (String) Azure Storage account name.
* `azure_blob_storage_settings.container_name` - (String) Name of the container in the Storage account.
* `azure_blob_storage_settingspath` - (String) Optional virtual sub-folder specifiying a path within the container.


