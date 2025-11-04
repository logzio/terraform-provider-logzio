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

variable "datasource_uid" {
  type        = string
  description = "Prometheus datasource UID in Logz.io"
}

variable "notification_email" {
  type        = string
  description = "Email address for alert notifications"
  default     = "alerts@example.com"
}

provider "logzio" {
  api_token = var.api_token
}

# Example 1: Basic threshold metric alert
resource "logzio_unified_alert" "high_cpu_alert" {
  title       = "High CPU Usage"
  type        = "METRIC_ALERT"
  description = "Fires when CPU usage exceeds 80%"
  tags        = ["infrastructure", "cpu"]
  enabled     = true

  metric_alert {
    severity = "HIGH"

    trigger {
      trigger_type             = "THRESHOLD"
      metric_operator          = "ABOVE"
      min_threshold            = 80.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "avg(rate(cpu_usage_total[5m])) * 100"
      }
    }

    recipients {
      emails = [var.notification_email]
    }
  }
}

# Example 2: Metric alert with range operator
resource "logzio_unified_alert" "memory_range_alert" {
  title       = "Memory Usage Outside Normal Range"
  type        = "METRIC_ALERT"
  description = "Alert when memory usage is outside 20-80% range"
  tags        = ["infrastructure", "memory"]
  enabled     = true

  metric_alert {
    severity = "MEDIUM"

    trigger {
      trigger_type             = "THRESHOLD"
      metric_operator          = "OUTSIDE_RANGE"
      min_threshold            = 20.0
      max_threshold            = 80.0
      search_timeframe_minutes = 10
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "avg(memory_usage_percent)"
      }
    }

    recipients {
      emails = [var.notification_email, "infrastructure@example.com"]
    }
  }
}

# Example 3: Math expression metric alert
resource "logzio_unified_alert" "error_rate_percent_alert" {
  title       = "5xx Error Rate Percentage High"
  type        = "METRIC_ALERT"
  description = "Fires when 5xx rate exceeds 2% of total requests"
  tags        = ["http", "errors", "api"]
  enabled     = true

  # RCA configuration
  rca = true
  use_alert_notification_endpoints_for_rca = true
  runbook = "Check API service health and recent deployments"

  metric_alert {
    severity = "HIGH"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "($A / $B) * 100"
      metric_operator          = "ABOVE"
      min_threshold            = 2.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "sum(rate(http_requests_total{status=~\"5..\"}[5m]))"
      }
    }

    queries {
      ref_id = "B"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "sum(rate(http_requests_total[5m]))"
      }
    }

    recipients {
      emails = [var.notification_email]
    }
  }
}

# Example 4: Complex math expression with multiple queries
resource "logzio_unified_alert" "latency_vs_baseline_alert" {
  title       = "Request Latency Above Baseline"
  type        = "METRIC_ALERT"
  description = "Alert when current latency exceeds baseline by 50%"
  tags        = ["performance", "latency"]
  enabled     = true

  metric_alert {
    severity = "MEDIUM"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "(($A - $B) / $B) * 100"
      metric_operator          = "ABOVE"
      min_threshold            = 50.0
      search_timeframe_minutes = 15
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
      }
    }

    queries {
      ref_id = "B"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "avg_over_time(http_request_duration_seconds[1h])"
      }
    }

    recipients {
      emails = [var.notification_email]
    }
  }
}

# Example 5: Low threshold alert
resource "logzio_unified_alert" "low_throughput_alert" {
  title       = "Request Throughput Below Minimum"
  type        = "METRIC_ALERT"
  description = "Alert when request rate drops below expected minimum"
  tags        = ["performance", "throughput"]
  enabled     = true

  metric_alert {
    severity = "LOW"

    trigger {
      trigger_type             = "THRESHOLD"
      metric_operator          = "BELOW"
      min_threshold            = 100.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "sum(rate(http_requests_total[5m]))"
      }
    }

    recipients {
      emails = [var.notification_email]
    }
  }
}

# Data source example - retrieve an existing alert
data "logzio_unified_alert" "existing_metric_alert" {
  type     = "METRIC_ALERT"
  alert_id = logzio_unified_alert.high_cpu_alert.alert_id
}

output "high_cpu_alert_id" {
  value       = logzio_unified_alert.high_cpu_alert.alert_id
  description = "ID of the high CPU alert"
}

output "error_rate_alert_id" {
  value       = logzio_unified_alert.error_rate_percent_alert.alert_id
  description = "ID of the error rate percentage alert"
}

output "existing_alert_severity" {
  value       = data.logzio_unified_alert.existing_metric_alert.metric_alert[0].severity
  description = "Severity of the retrieved alert"
}

