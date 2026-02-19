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

# Example 1: Basic log alert with COUNT aggregation and JSON output
resource "logzio_unified_alert" "basic_log_alert" {
  title       = "High Error Rate"
  description = "Triggers when error logs exceed threshold"
  tags        = ["production", "errors"]
  enabled     = true

  recipients {
    emails = [var.notification_email]
  }

  alert_configuration {
    type                           = "LOG_ALERT"
    search_timeframe_minutes       = 15
    suppress_notifications_minutes = 30
    alert_output_template_type     = "JSON"

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

# Example 2: Log alert with TABLE output and custom columns
# TABLE output requires aggregation_type = "NONE", should_use_all_fields = false,
# and at least one column defined.
resource "logzio_unified_alert" "table_log_alert" {
  title       = "Critical Service Errors"
  description = "Alert on critical errors in payment service with table output"
  tags        = ["critical", "payment", "service"]
  enabled     = true

  # RCA configuration
  rca                                      = true
  use_alert_notification_endpoints_for_rca = true
  runbook = "1. Check payment gateway status\n2. Review recent deployments\n3. Contact on-call engineer"

  recipients {
    emails = [var.notification_email, "oncall@example.com"]
  }

  alert_configuration {
    type                           = "LOG_ALERT"
    search_timeframe_minutes       = 10
    suppress_notifications_minutes = 60
    alert_output_template_type     = "TABLE"

    sub_components {
      query_definition {
        query                        = "service:payment AND level:CRITICAL"
        should_query_on_all_accounts = false
        account_ids_to_query_on      = [12345]

        aggregation {
          aggregation_type = "NONE"
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
          field_name = "@timestamp"
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

    # Custom schedule - check every 5 minutes
    schedule {
      cron_expression = "0 0/5 * * * ?"
      timezone        = "America/New_York"
    }
  }
}

# Example 3: Log alert with multiple sub-components and correlation
resource "logzio_unified_alert" "correlated_log_alert" {
  title       = "Database Connection Issues"
  description = "Alert when both connection errors and timeouts occur"
  tags        = ["database", "connections"]
  enabled     = true

  recipients {
    emails = [var.notification_email]
  }

  alert_configuration {
    type                       = "LOG_ALERT"
    search_timeframe_minutes   = 5
    alert_output_template_type = "JSON"

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
