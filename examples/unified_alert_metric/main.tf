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

variable "metrics_account_id" {
  type        = number
  description = "Logz.io Metrics account ID"
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
  description = "Fires when CPU usage exceeds 80%"
  tags        = ["infrastructure", "cpu"]
  enabled     = true

  recipients {
    emails = [var.notification_email]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "HIGH"

    trigger {
      type = "threshold"

      condition {
        operator_type = "above"
        threshold     = 80.0
      }
    }

    queries {
      ref_id = "A"

      query_definition {
        account_id   = var.metrics_account_id
        promql_query = "avg(rate(cpu_usage_total[5m])) * 100"
      }
    }
  }
}

# Example 2: Metric alert with range operator
resource "logzio_unified_alert" "memory_range_alert" {
  title       = "Memory Usage Outside Normal Range"
  description = "Alert when memory usage is outside 20-80% range"
  tags        = ["infrastructure", "memory"]
  enabled     = true

  recipients {
    emails = [var.notification_email, "infrastructure@example.com"]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "MEDIUM"

    trigger {
      type = "threshold"

      condition {
        operator_type = "outside_range"
        from          = 20.0
        to            = 80.0
      }
    }

    queries {
      ref_id = "A"

      query_definition {
        account_id   = var.metrics_account_id
        promql_query = "avg(memory_usage_percent)"
      }
    }
  }
}

# Example 3: Math expression metric alert
resource "logzio_unified_alert" "error_rate_percent_alert" {
  title       = "5xx Error Rate Percentage High"
  description = "Fires when 5xx rate exceeds 2% of total requests"
  tags        = ["http", "errors", "api"]
  enabled     = true

  # RCA configuration
  rca                                      = true
  use_alert_notification_endpoints_for_rca = true
  runbook = "Check API service health and recent deployments"

  recipients {
    emails = [var.notification_email]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "HIGH"

    trigger {
      type       = "math"
      expression = "($A / $B) * 100"
    }

    queries {
      ref_id = "A"

      query_definition {
        account_id   = var.metrics_account_id
        promql_query = "sum(rate(http_requests_total{status=~\"5..\"}[5m]))"
      }
    }

    queries {
      ref_id = "B"

      query_definition {
        account_id   = var.metrics_account_id
        promql_query = "sum(rate(http_requests_total[5m]))"
      }
    }
  }
}

# Example 4: Complex math expression with multiple queries
resource "logzio_unified_alert" "latency_vs_baseline_alert" {
  title       = "Request Latency Above Baseline"
  description = "Alert when current latency exceeds baseline by 50%"
  tags        = ["performance", "latency"]
  enabled     = true

  recipients {
    emails = [var.notification_email]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "MEDIUM"

    trigger {
      type       = "math"
      expression = "(($A - $B) / $B) * 100"
    }

    queries {
      ref_id = "A"

      query_definition {
        account_id   = var.metrics_account_id
        promql_query = "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
      }
    }

    queries {
      ref_id = "B"

      query_definition {
        account_id   = var.metrics_account_id
        promql_query = "avg_over_time(http_request_duration_seconds[1h])"
      }
    }
  }
}

# Example 5: Low threshold alert
resource "logzio_unified_alert" "low_throughput_alert" {
  title       = "Request Throughput Below Minimum"
  description = "Alert when request rate drops below expected minimum"
  tags        = ["performance", "throughput"]
  enabled     = true

  recipients {
    emails = [var.notification_email]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "LOW"

    trigger {
      type = "threshold"

      condition {
        operator_type = "below"
        threshold     = 100.0
      }
    }

    queries {
      ref_id = "A"

      query_definition {
        account_id   = var.metrics_account_id
        promql_query = "sum(rate(http_requests_total[5m]))"
      }
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
  value       = data.logzio_unified_alert.existing_metric_alert.alert_configuration[0].severity
  description = "Severity of the retrieved alert"
}
