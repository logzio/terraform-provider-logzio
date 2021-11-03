resource "logzio_alert_v2" "%s" {
  title = "updated_alert"
  description = "this is a description"
  tags = [
    "some",
    "test"]
  search_timeframe_minutes = 5
  is_enabled = false
  notification_emails = [
    "testx@test.com"]
  suppress_notifications_minutes = 5
  output_type = "JSON"
  sub_components {
    query_string = "loglevel:ERROR"
    should_query_on_all_accounts = true
    operation = "EQUALS"
    value_aggregation_type = "COUNT"
    severity_threshold_tiers {
      severity = "HIGH"
      threshold = 10
    }
    severity_threshold_tiers {
      severity = "LOW"
      threshold = 2
    }
    filter_must = jsonencode([
      {
        match_phrase: {
          "some_field": {
            "query": "some_query"
          }
        }
      },
      {
        match_phrase: {
          "some_field2": {
            "query": "hello world"
          }
        }
      }
    ])
  }
}