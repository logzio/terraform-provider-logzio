# Alert Provider

Provides a Logz.io log monitoring alert resource. This can be used to create and manage Logz.io log monitoring alerts. 

* Learn more about log alerts in the [Logz.io Docs](https://docs.logz.io/user-guide/alerts/)

## Example Usage

```hcl
# Create a new alert and a new endpoint
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

*	`alert_notification_endpoints` - (Required) 
*	`description` - (Optional) A description of the event, its significance, and suggested next steps or instructions for the team. 
*	`filter` - (Optional) 
*	`tags` - (Optional) Tags for filtering alerts and triggered alerts. Can be used in Kibana Discover, Kibana dashboards, and more.
*	`group_by_aggregation_fields` - (Optional)
*	`is_enabled` - (Optional) True by default. If `true`, the alert is currently active.
*	`query_string` - (Optional)
*	`last_triggered_at` - (Optional)
*	`last_updated` - (Optional) Date and time in UTC when the alert was last updated. | 
*	`notification_emails` - (Optional) 
*	`operation` - (Optional) 
*	`search_timeframe_minutes` - (Optional)  The time frame for evaluating the log data is a sliding window, with 1 minute granularity.

  The recommended minimum and maximum values are not validated, but needed to guarantee the alert's accuracy.

  The minimum recommended time frame is 5 minutes, as anything shorter will be less reliable and unnecessarily resource-heavy.

  The maximum recommended time frame is 1440 minutes (24 hours). The alert runs on the index from today and yesterday (in UTC) and the maximum time frame increases throughout the day, reaching 48 hours exactly before midnight UTC.  
*	`severity` - (Optional)
*	`severity_threshold_tiers` - (Optional)
*	`suppress_notifications_minutes` - (Optional)
*	`threshold` - (Optional)
*	`title` - (Optional)
* `Alert title` - (Optional) 
*	`value_aggregation_field` - (Optional)
* `value_aggregation_type` - (Optional)

## Attribute Reference

*	`id` - (Required) Logz.io alert ID. 
*	`created_at` - (Required) Date and time in UTC when the alert was first created.
*	`created_by` - (Optional) Email of the user who first created the alert.
