variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = "${var.api_token}"
}

resource "logzio_endpoint" "my_endpoint" {
  title = "my_endpoint"
  description = "hello"
  endpoint_type = "slack"
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
  alert_notification_endpoints = ["${logzio_endpoint.my_endpoint.id}"]
  suppress_notifications_minutes = 5
  severity_threshold_tiers = [
    {
      "severity" = "HIGH",
      "threshold" = 10
    }
  ]
}