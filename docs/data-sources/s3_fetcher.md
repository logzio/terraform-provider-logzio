# S3 Fetcher Datasource

Use this data source to access information about existing Logz.io S3 Fetchers.

## Argument Reference

* `fetcher_id` - ID of the S3 Fetcher.

##  Attribute Reference

* `aws_access_key` - (String) AWS S3 bucket access key. Not applicable if you chose to authenticate with `aws_arn`.
* `aws_arn` - (String) Amazon Resource Name (ARN) to uniquely identify the S3 bucket. Not applicable if you choose to authenticate with AWS keys (access key & secret key).
* `bucket_name` - (String) AWS S3 bucket name.
* `active` - (Boolean) If true, the S3 bucket connector is active and logs are being fetched to Logz.io. If false, the connector is disabled.
* `aws_region` - (String) Bucket's region. Allowed values: `US_EAST_1`, `US_EAST_2`, `US_WEST_1`, `US_WEST_2`, `EU_WEST_1`, `EU_WEST_2`, `EU_WEST_3`, `EU_CENTRAL_1`, `AP_NORTHEAST_1`, `AP_NORTHEAST_2`, `AP_SOUTHEAST_1`, `AP_SOUTHEAST_2`, `SA_EAST_1`, `AP_SOUTH_1`, `CA_CENTRAL_1`.
* `logs_type` - (String) Specifies the log type being sent to Logz.io. Determines the parsing pipeline used to parse and map the logs. [Learn more about parsing options supported by Logz.io](https://docs.logz.io/user-guide/log-shipping/built-in-log-types.html). Allowed values: `elb`, `vpcflow`, `S3Access`, `cloudfront`.
* `prefix` - (String) Prefix of the AWS S3 bucket.
* `add_s3_object_key_as_log_field` - (Boolean) If `true`, enriches logs with a new field detailing the S3 object key.
