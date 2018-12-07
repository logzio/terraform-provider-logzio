# Logz.io terraform provider

### Supports only alerts at the moment!

##### Using the provider

A simple example, will create an alert, triggered if over 5 minutes more than 10 logs are found with a loglevel of ERROR,
further notifications are suppressed for 5 minutes.

You can also provide the `api_token` as the `LOGZIO_API_TOKEN` env var.

```hcl-terraform
provider "logzio" {
  api_token = "${var.api_token}"
}

resource "logzio_alert" "my_alert" {
  title = "my_alert_title"
  query_string = "loglevel:ERROR"
  operation = "GREATER_THAN"
  notification_emails = ["${var.notification_email}"]
  search_timeframe_minutes = 5
  value_aggregation_type = "NONE"
  alert_notification_endpoints = []
  suppress_notifications_minutes = 5
  severity_threshold_tiers = [
    {
      "severity" = "HIGH",
      "threshold" = 10
    }
   ]
}
```



##### Logz.io API support

|api  |method|support     |implementation|
|-----|------|------------|--------------|
|alert|create|implemented |`resource_alert::resourceAlertCreate`|
|alert|delete|implemented |`resource_alert::resourceAlertDelete`|
|alert|update|implemented |`resource_alert::resourceAlertUpdate`|
|alert|read  |implemented |`resource_alert::resourceAlertRead`  |

##### Building and testing the provider

Run `scripts/build.sh` to build the provider, copy it to the terraform execution directory (the `terraform` directory) and
do a complete terraform init/plan/apply lifecycle.

Run `scripts/destroy.sh` to build the provider, copy it and run the destroy part of the terraform lifecycle.

##### Terraform demos

The terraform demos are in the `terraform` directory.

To use the demos you'll need to provide an `api_token` and a `notification_email` as variables.

There's a utility script in the `terraform` directory that will run the terraform init/plan/apply lifecycle the templates in that directory.  The script looks for a file called `variables\${USER}.tfvars` to set the values of `api_token` and `notification_email`.

