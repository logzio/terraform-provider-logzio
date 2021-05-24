resource "logzio_alert_v2" "%s" {
  title = "updated_alert"
  description = "this is a description"
  tags = ["some", "test"]
  search_timeframe_minutes = 5
  is_enabled = false
  notification_emails = ["testx@test.com"]
  suppress_notifications_minutes = 5
  output_type = "JSON"
  sub_components {
    query_string = "loglevel:ERROR"
    should_query_on_all_accounts = true
    operation = "GREATER_THAN"
    value_aggregation_type = "SUM"
    value_aggregation_field = "some_field"
    severity_threshold_tiers {
      severity = "HIGH"
      threshold = 10
    }
  }
}