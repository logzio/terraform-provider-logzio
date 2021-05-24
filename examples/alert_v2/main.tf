variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_alert_v2" "my_alert" {
  title = "hello_there"
  search_timeframe_minutes = 5
  is_enabled = false
  tags = ["some", "words"]
  suppress_notifications_minutes = 5
  output_type = "JSON"
  sub_components {
    query_string = "loglevel:ERROR"
    should_query_on_all_accounts = true
    operation = "GREATER_THAN"
    value_aggregation_type = "COUNT"
    severity_threshold_tiers {
      severity = "HIGH"
      threshold = 10
    }
    severity_threshold_tiers {
      severity = "INFO"
      threshold = 5
    }
  }
}