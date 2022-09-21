# Archive Logs Provider

Provides a Logz.io archive logs resource. This can be used to create and manage Logz.io "archive logs" settings.

* Learn more about archive logs in the [Logz.io Docs](https://docs.logz.io/api/#tag/Archive-logs)

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

```

## Argument Reference

### Required:

* `storage_type` - (String) Specifies the storage provider. Applicable values: `S3`, `BLOB`.

### Optional:

* `enabled` - (Boolean) Defaults to `true`. If `true`, archiving is currently enabled.
* `compressed` - (Boolean) Defaults to `true`. If `true`, logs are compressed before they are archived.

#### Required if `storage_type` is `S3`:

* `credentials_type` - (String) Specifies which credentials will be used for authentication.
The options are either `KEYS` or `IAM`. Authentication with S3 Secret Credentials is supported for backward compatibility. IAM roles are strongly recommended.
* `s3_path` - (String) Specify a path to the **root** of the S3 bucket. (Currently, archiving to a sub-bucket is supported, but not restoring from one.) **Unique buckets** - It is important to archive each account/sub-account to a separate S3 bucket.
* `s3_iam_credentials_arn` - (String) Applicable when `credentials_type` is `IAM`. Amazon Resource Name (ARN) to uniquely identify the S3 bucket.
* `aws_access_key` - (String) Applicable when `credentials_type` is `KEYS`.
* `aws_secret_key` - (String) Applicable when `credentials_type` is `KEYS`.

##### Required if `storage_type` is `BLOB`:

* `tenant_id` - (String) Azure Directory (tenant) ID. The Tenant ID of the AD app. Go to **Azure Active Directory > App registrations** and select the app to see it.
* `client_id` - (String) Azure application (client) ID. The Client ID of the AD app, found under the App Overview page. Go to **Azure Active Directory > App registrations** and select the app to see it.
* `client_secret` - (String) Azure client secret. Password of the Client secret, found in the app's **Certificates & secrets** page. Go to **Azure Active Directory > App registrations** and select the app. Then select **Certificates & secrets** to see it.
* `account_name` - (String) Azure Storage account name. Name of the storage account that holds the container where the logs will be archived.
* `container_name` - (String) Name of the container in the Storage account. This is where the logs will be archived.

##### Optional:

* `path` - (String) Optional virtual sub-folder specifiying a path within the container. Logs will be archived under the path “{container-name}/{virtual sub-folder}”. Avoid leading and trailing slashes (/). For example, the prefix “region1” is good, but “/region1/” is not.

##  Attribute Reference

* `archive_id` - (String) Archive ID in the Logz.io database.

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
