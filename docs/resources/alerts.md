# Alert Provider

Provides a Logz.io log monitoring alert resource. This can be used to create and manage Logz.io log monitoring alerts.

* Learn more about log alerts in the [Logz.io Docs](https://docs.logz.io/user-guide/alerts/)

## Example Usage

```hcl
# Create a new alert and a new endpoint for use as the alert notification channel
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_endpoint" "my_endpoint" {
  title = "my_endpoint"
  description = "hello"
  endpoint_type = "Slack"
  slack {
    url = "https://this.is.com/some/url"
  }
}

resource "logzio_alert" "my_alert" {
  title = "my_other_title"
  query_string = "loglevel:ERROR"
  operation = "GREATER_THAN"
  notification_emails = []
  search_timeframe_minutes = 5
  value_aggregation_type = "NONE"
  alert_notification_endpoints = [logzio_endpoint.my_endpoint.id]
  suppress_notifications_minutes = 5
  is_enabled = false
  severity_threshold_tiers {
      severity = "HIGH"
      threshold = 100
    }
  severity_threshold_tiers {
    severity = "LOW"
    threshold = 20
  }
  tags = ["some", "words"]
}

```

## Argument Reference

*	`title` - (Required) Alert title.
*	`search_timeframe_minutes` - (Required)  The time frame for evaluating the log data is a sliding window, with 1 minute granularity.
*	`operation` - (Optional) Specifies the operator for evaluating the results. Enum: `LESS_THAN`, `GREATER_THAN`, `LESS_THAN_OR_EQUALS`, `GREATER_THAN_OR_EQUALS`, `EQUALS`, `NOT_EQUALS`.
*	`severity_threshold_tiers` - (Required) Set as many as
5 thresholds, each with its own severity level. Map schema sets a severity per trigger threshold. 

  If setting several thresholds, you must adhere to the logic of the operator. See the [API docs](https://docs.logz.io/api/#operation/createAlert) for futher details.
  *	`severity` - Defaults to `MEDIUM`. Labels the event with a severity tag. Available severity tags are: `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
  *	`threshold` - (Required) 

*	`alert_notification_endpoints` - (Required) Add email addresses and/or endpoint channels to automatically receive notifications with sample data when the alert triggers.
*	`description` - (Optional) A description of the event, its significance, and suggested next steps or instructions for the team.
*	`query_string` - (Required) Search query in Lucene syntax. You can combine filters and a search query to specify the logs you are looking for. You can combine filters and a search query to specify the logs you are looking for.
*	`filter` - (Optional) You can use `must` and `must_not` filters. Filters are more efficient compared to a query, so it's recommended to opt for a filter over a `query_string`, where possible.
*	`tags` - (Optional) Tags for filtering alerts and triggered alerts. Can be used in Kibana Discover, Kibana dashboards, and more.
*	`group_by_aggregation_fields` - (Optional) Specify 1-3 fields by which to group the results and count them. If you apply a group by operation, the alert returns a count of the results aggregated by unique values.
*	`is_enabled` - (Optional) True by default. If `true`, the alert is currently active.
*	`last_triggered_at` - (Optional) Date and time in UTC when the alert last triggered.
*	`last_updated` - (Optional) Date and time in UTC when the alert was last updated.
*	`notification_emails` - (Optional) Array of email addresses to be notified when the alert triggers.
*	`suppress_notifications_minutes` - (Optional)
*	`value_aggregation_field` - (Optional) Specify the field on which to run the aggregation for the trigger condition.
* `value_aggregation_type` - (Required) Specifies the aggregation operator. Can be: `SUM`, `MIN`, `MAX`, `AVG`, `COUNT`, `UNIQUE_COUNT`, `NONE`.

## Attribute Reference

*	`id` - Logz.io alert ID.
*	`created_at` - Date and time in UTC when the alert was first created.
*	`created_by` - Email of the user who first created the alert.
