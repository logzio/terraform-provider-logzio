resource "logzio_alert_v2" "alert_v2_datasource" {
  title = "hello"
  search_timeframe_minutes = 5
  is_enabled = false
  notification_emails = ["testx@test.com"]
  suppress_notifications_minutes = 5
  output_type = "JSON"
  sub_components {
    query_string = "loglevel:ERROR"
    operation = "GREATER_THAN"
    value_aggregation_type = "NONE"
    severity_threshold_tiers {
      severity = "HIGH"
      threshold = 10
    }
  }
}

data "logzio_alert_v2" "alert_v2_datasource_by_id" {
  id = "${logzio_alert_v2.alert_v2_datasource.id}"
  depends_on = ["logzio_alert_v2.alert_v2_datasource"]
}