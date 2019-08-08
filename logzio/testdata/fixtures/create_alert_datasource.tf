resource "logzio_alert" "alert_datasource" {
  title = "hello"
  query_string = "loglevel:ERROR"
  operation = "GREATER_THAN"
  notification_emails = ["testx@test.com"]
  search_timeframe_minutes = 5
  value_aggregation_type = "NONE"
  alert_notification_endpoints = []
  suppress_notifications_minutes = 5
  severity_threshold_tiers {
    severity = "HIGH"
    threshold = 10
  }
}

data "logzio_alert" "alert_datasource_by_id" {
  id = "${logzio_alert.alert_datasource.id}"
  depends_on = ["logzio_alert.alert_datasource"]
}