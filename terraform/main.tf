provider "logzio" {
  api_token = "${var.api_token}"
}

resource "logzio_alert" "my_alert" {
  title = "my_other_title"
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