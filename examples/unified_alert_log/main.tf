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
  description = "Your Logz.io API token"
}

variable "notification_email" {
  type        = string
  description = "Email address for alert notifications"
  default     = "alerts@example.com"
}

provider "logzio" {
  api_token = var.api_token
}

# Example 1: Basic log alert with COUNT aggregation
resource "logzio_unified_alert" "basic_log_alert" {
  title       = "High Error Rate"
  type        = "LOG_ALERT"
  description = "Triggers when error logs exceed threshold"
  tags        = ["production", "errors"]
  enabled     = true

  log_alert {
    search_timeframe_minutes = 15

    output {
      type                          = "JSON"
      suppress_notifications_minutes = 30

      recipients {
        emails = [var.notification_email]
      }
    }

    sub_components {
      query_definition {
        query                        = "level:ERROR"
        should_query_on_all_accounts = true

        aggregation {
          aggregation_type = "COUNT"
        }
      }

      trigger {
        operator = "GREATER_THAN"

        severity_threshold_tiers {
          severity  = "HIGH"
          threshold = 100
        }

        severity_threshold_tiers {
          severity  = "MEDIUM"
          threshold = 50
        }

        severity_threshold_tiers {
          severity  = "LOW"
          threshold = 10
        }
      }

      output {
        should_use_all_fields = true
      }
    }
  }
}

# Example 2: Advanced log alert with specific accounts and custom columns
resource "logzio_unified_alert" "advanced_log_alert" {
  title       = "Critical Service Errors"
  type        = "LOG_ALERT"
  description = "Alert on critical errors in payment service"
  tags        = ["critical", "payment", "service"]
  enabled     = true

  # RCA configuration
  rca = true
  use_alert_notification_endpoints_for_rca = true
  runbook = "1. Check payment gateway status\n2. Review recent deployments\n3. Contact on-call engineer"

  log_alert {
    search_timeframe_minutes = 10

    output {
      type                          = "TABLE"
      suppress_notifications_minutes = 60

      recipients {
        emails = [var.notification_email, "oncall@example.com"]
      }
    }

    sub_components {
      query_definition {
        query                        = "service:payment AND level:CRITICAL"
        should_query_on_all_accounts = false
        account_ids_to_query_on      = [12345]

        aggregation {
          aggregation_type    = "SUM"
          field_to_aggregate_on = "error_count"
        }

        group_by = ["service", "error_type"]
      }

      trigger {
        operator = "GREATER_THAN_OR_EQUALS"

        severity_threshold_tiers {
          severity  = "SEVERE"
          threshold = 50
        }

        severity_threshold_tiers {
          severity  = "HIGH"
          threshold = 20
        }
      }

      output {
        should_use_all_fields = false

        columns {
          field_name = "timestamp"
          sort       = "DESC"
        }

        columns {
          field_name = "service"
        }

        columns {
          field_name = "error_type"
        }

        columns {
          field_name = "message"
        }
      }
    }

    # Custom schedule - check every 5 minutes during business hours
    schedule {
      cron_expression = "*/5 9-17 * * 1-5"
      timezone        = "America/New_York"
    }
  }
}

# Example 3: Log alert with multiple sub-components and correlation
resource "logzio_unified_alert" "correlated_log_alert" {
  title       = "Database Connection Issues"
  type        = "LOG_ALERT"
  description = "Alert when both connection errors and timeouts occur"
  tags        = ["database", "connections"]
  enabled     = true

  log_alert {
    search_timeframe_minutes = 5

    output {
      type = "JSON"

      recipients {
        emails = [var.notification_email]
      }
    }

    sub_components {
      query_definition {
        query                        = "database:postgres AND error:connection"
        should_query_on_all_accounts = true

        aggregation {
          aggregation_type = "COUNT"
        }
      }

      trigger {
        operator = "GREATER_THAN"

        severity_threshold_tiers {
          severity  = "HIGH"
          threshold = 10
        }
      }

      output {
        should_use_all_fields = true
      }
    }

    sub_components {
      query_definition {
        query                        = "database:postgres AND timeout:true"
        should_query_on_all_accounts = true

        aggregation {
          aggregation_type = "COUNT"
        }
      }

      trigger {
        operator = "GREATER_THAN"

        severity_threshold_tiers {
          severity  = "HIGH"
          threshold = 5
        }
      }

      output {
        should_use_all_fields = true
      }
    }

    correlations {
      correlation_operators = ["AND"]
    }
  }
}

# Data source example - retrieve an existing alert
data "logzio_unified_alert" "existing_alert" {
  type     = "LOG_ALERT"
  alert_id = logzio_unified_alert.basic_log_alert.alert_id
}

output "basic_alert_id" {
  value       = logzio_unified_alert.basic_log_alert.alert_id
  description = "ID of the basic log alert"
}

output "basic_alert_created_at" {
  value       = logzio_unified_alert.basic_log_alert.created_at
  description = "Creation timestamp of the basic log alert"
}

output "existing_alert_title" {
  value       = data.logzio_unified_alert.existing_alert.title
  description = "Title of the retrieved alert"
}

