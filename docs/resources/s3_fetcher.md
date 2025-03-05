# S3 Fetcher Provider

Provides a Logz.io S3 Fetcher resource. This can be used to create and manage S3 Fetcher instances, that retrieve logs stored in AWS S3 Buckets.

* Learn more about S3 Fetcher in the [Logz.io Docs](https://docs.logz.io/api/#tag/Connect-to-S3-Buckets).

## Example Usage

```hcl
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
```

## Argument Reference

### Required:

* `aws_access_key` - (String) AWS S3 bucket access key. Not applicable if you choose to authenticate with `aws_arn`. If you choose to authenticate with AWS keys, both `aws_access_key` and `aws_secret key` must be set.
* `aws_secret_key` - (String) AWS S3 bucket secret key. Not applicable if you choose to authenticate with `aws_arn`. If you choose to authenticate with AWS keys, both `aws_access_key` and `aws_secret key` must be set.
* `aws_arn` - (String) Amazon Resource Name (ARN) of the IAM Role used for authentication. To generate a new ARN, create a new IAM Role in your AWS admin console. Not applicable if you choose to authenticate with AWS keys (access key & secret key).
* `bucket_name` - (String) AWS S3 bucket name.
* `active` - (Boolean) If true, the S3 bucket connector is active and logs are being fetched to Logz.io. If false, the connector is disabled.
* `aws_region` - (String) Bucket's region. Specify one supported AWS region: `US_EAST_1`, `US_EAST_2`, `US_WEST_1`, `US_WEST_2`, `EU_WEST_1`, `EU_WEST_2`, `EU_WEST_3`, `EU_CENTRAL_1`, `AP_NORTHEAST_1`, `AP_NORTHEAST_2`, `AP_SOUTHEAST_1`, `AP_SOUTHEAST_2`, `SA_EAST_1`, `AP_SOUTH_1`, `CA_CENTRAL_1`.
* `logs_type` - (String) Specifies the log type being sent to Logz.io. Determines the parsing pipeline used to parse and map the logs. [Learn more about parsing options supported by Logz.io](https://docs.logz.io/user-guide/log-shipping/built-in-log-types.html). Allowed values: `elb`, `vpcflow`, `S3Access`, `cloudfront`.

### Optional:

* `prefix` - (String) Prefix of the AWS S3 bucket.
* `add_s3_object_key_as_log_field` - (Boolean) Defaults to `false`. If `true`, enriches logs with a new field detailing the S3 object key.


## Attribute Reference

* `fetcher_id` - (Int) Id of the created fetcher.

### Import S3 Fetcher as Terraform resource

You can import existing S3 Fetcher as follows:

```
terraform import logzio_s3_fetcher.my_fetcher <FETCHER-ID>
```