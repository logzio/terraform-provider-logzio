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
*	`alert_notification_endpoints` - (Optional) 
*	`description` - (Optional) A description of the event, its significance, and suggested next steps or instructions for the team. 
*	`filter` - (Optional) 
*	`tags` - (Optional) Tags for filtering alerts and triggered alerts. Can be used in Kibana Discover, Kibana dashboards, and more.
*	`group_by_aggregation_fields` - (Optional)
*	`is_enabled` - (Optional) True by default. If `true`, the alert is currently active.
*	`query_string` - (Optional) Search query in Lucene syntax. Determines when the alert should trigger in combination with filters, group by aggregations, accounts to search, and trigger conditions.
*	`last_triggered_at` - (Optional) Date and time in UTC when the alert last triggered.
*	`last_updated` - (Optional) Date and time in UTC when the alert was last updated.
*	`notification_emails` - (Optional) Array of email addresses to be notified when the alert triggers.
*	`operation` - (Optional) Specifies the operator for evaluating the results. Enum: `LESS_THAN`, `GREATER_THAN`, `LESS_THAN_OR_EQUALS`, `GREATER_THAN_OR_EQUALS`, `EQUALS`, `NOT_EQUALS`.
*	`search_timeframe_minutes` - (Required)  The time frame for evaluating the log data is a sliding window, with 1 minute granularity.
*	`severity` - (Optional) Defaults to `MEDIUM`. Specifies a severity for the event when the alert triggers. Can be `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
*	`threshold` - (Optional) 
*	`severity_threshold_tiers` - (Optional) Set per trigger threshold.
*	`suppress_notifications_minutes` - (Optional) 
*	`value_aggregation_field` - (Optional) 
* `value_aggregation_type` - (Required)

## Attribute Reference

*	`id` - Logz.io alert ID. 
*	`created_at` - Date and time in UTC when the alert was first created.
*	`created_by` - Email of the user who first created the alert.
