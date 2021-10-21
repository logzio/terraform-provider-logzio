# Archive Logs Provider

Provides a Logz.io archive logs resource. This can be used to create and manage Logz.io archive logs settings.

* Learn more about log shipping tokens in the [Logz.io Docs](https://docs.logz.io/api/#tag/Archive-logs)

## Example Usage

```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_archive_logs" "my_s3_archive" { 
  storage_type = "S3"
  compressed = true
  enabled = true
  amazon_s3_storage_settings { 
    credentials_type = "IAM"
    s3_path = "my-s3-path"
    s3_iam_credentials_arn = "my-arn"
 }
}

```

## Argument Reference

### Required:
* `storage_type` - (String) Specifies the storage provider. Applicable values: `S3`, `BLOB`.
If `S3`, the `amazon_s3_storage_settings` are relevant. 
If `BLOB`, the `azure_blob_storage_settings` are relevant.
* `amazon_s3_storage_settings` - (Object) Applicable settings when the `storage_type` is `S3`. See below for **nested schema**.
* `azure_blob_storage_settings` - (Object) Applicable settings when the `storage_type` is `BLOB`. See below for **nested schema**.

### Optional:
* `enabled` - (Boolean) Defaults to `true`. If `true`, archiving is currently enabled.
* `compressed` - (Boolean) Defaults to `true`. If `true`, logs are compressed before they are archived.

#### Nested schema for `amazon_s3_storage_settings`:

* `credentials_type` - (String) Specifies which credentials will be used for authentication.
The options are either `KEYS` with `s3_secret_credentials`, or `IAM` with `s3_iam_credentials_arn`.
* `path` - (String) Specify a path to the **root** of the S3 bucket. (Currently, archiving to a sub-bucket is supported, but not restoring from one.) **Unique buckets** - It is important to archive each account/sub-account to a separate S3 bucket.
* `s3_secret_credentials` - (Object) Applicable settings when the `credentials_type` is `KEYS`.
Authentication with S3 Secret Credentials is supported for backward compatibility.
IAM roles are strongly recommended. See below for **nested schema**.
* `s3_iam_credentials_arn` - (String) Amazon Resource Name (ARN) to uniquely identify the S3 bucket.

##### Nested schema for `s3_secret_credentials`:
* `access_key` - (String)
* `secret_key` - (String)

#### Nested schema for `azure_blob_storage_settings`:

##### Required:
* `tenant_id` - (String) Azure Directory (tenant) ID. The Tenant ID of the AD app. Go to Azure Active Directory > App registrations and select the app to see it.
* `client_id` - (String) Azure application (client) ID. The Client ID of the AD app, found under the App Overview page. Go to Azure Active Directory > App registrations and select the app to see it.
* `client_secret` - (String) Azure client secret. Password of the Client secret, found in the app's Certificates & secrets page. Go to Azure Active Directory > App registrations and select the app. Then select Certificates & secrets to see it.
* `account_name` - (String) Azure Storage account name. Name of the storage account that holds the container where the logs will be archived.
* `container_name` - (String) Name of the container in the Storage account. This is where the logs will be archived.

##### Optional:
* `path` - (String) Optional virtual sub-folder specifiying a path within the container. Logs will be archived under the path “{container-name}/{virtual sub-folder}”. Avoid leading and trailing slashes (/). For example, the prefix “region1” is good, but “/region1/” is not.

##  Attribute Reference
* `archive_id` - (String) Archive ID in the Logz.io database.
* `amazon_s3_storage_settings.s3_external_id` - (String) The external id that gives Logz.io access to your S3 bucket.

## Importing resource:
To import an archive you'll need to specify your archive's id.
For example, if you have in your Terraform configuration the following:

```hcl
resource "logzio_archive_logs" "imported" {
  ...
}
```

And your archives's id is 123456, then your import command should be:

```bash
terraform import logzio_archive_logs.imported 123456
```