terraform {
  required_providers {
    logzio = {
      source  = "logzio/logzio"
      version = "~> 1.0"
    }
  }
}

variable "api_token" {
  type        = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_alert_v2" "my_alert" {
  title                     = "hello"
  description               = "this is a description"
  tags                      = ["some", "test"]
  search_timeframe_minutes  = 5
  is_enabled                = false
  notification_emails       = ["testx@test.com"]
  suppress_notifications_minutes = 5
  output_type               = "JSON"
  sub_components {
    query_string              = "loglevel:ERROR"
    should_query_on_all_accounts = true
    operation                = "EQUALS"
    value_aggregation_type   = "COUNT"
    severity_threshold_tiers {
      severity  = "HIGH"
      threshold = 10
    }
    severity_threshold_tiers {
      severity  = "LOW"
      threshold = 2
    }
    filter_must = jsonencode([
      {
        "bool": {
          "should": [
            {
              "match_phrase": {
                "query": "some_query"
              }
            }
          ]
        }
      }
    ])
  }

  schedule_cron_expression = "0 0/5 9-17 ? * * *"
  schedule_timezone        = "IET"
}
