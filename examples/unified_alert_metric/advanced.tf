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
  sensitive   = true
}

variable "notification_email" {
  type        = string
  description = "Email address for alert notifications"
  default     = "test@example.com"
}

variable "datasource_uid" {
  type        = string
  description = "Prometheus datasource UID in Logz.io"
  default     = "prometheus"
}

variable "notification_endpoint_id" {
  type        = number
  description = "Notification endpoint ID for alerts"
  default     = null
}

variable "rca_notification_endpoint_id" {
  type        = number
  description = "Separate notification endpoint ID for RCA"
  default     = null
}

variable "dashboard_id" {
  type        = string
  description = "Dashboard UID for alert context"
  default     = ""
}
variable "panel_id" {
  type        = string
  description = "panel UID for alert context"
  default     = ""
}

variable "folder_id" {
  type        = string
  description = "Folder UID for alert organization"
  default     = ""
}

provider "logzio" {
  api_token = var.api_token
}

# Basic log alert example for testing
resource "logzio_unified_alert" "test_log_alert" {
  title       = "Test High Error Rate"
  type        = "LOG_ALERT"
  description = "Test alert for local provider testing"
  tags        = ["test", "local"]
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

# =============================================================================
# PRIORITIZED METRIC ALERT TEST CASES
# =============================================================================

# Case 1: One Query, No AI
# Simple threshold alert with single query
resource "logzio_unified_alert" "metric_one_query_no_ai" {
  title       = "Test: One Query No AI - High CPU"
  type        = "METRIC_ALERT"
  description = "Alert when CPU usage exceeds threshold - single query, no RCA"
  tags        = ["test", "metrics", "cpu", "no-ai"]
  enabled     = true

  # Optional dashboard linking
  dashboard_id = var.dashboard_id != "" ? var.dashboard_id : null
  folder_id    = var.folder_id != "" ? var.folder_id : null
  panel_id = var.folder_id

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
        promql_query   = "avg(rate(node_cpu_seconds_total{mode=\"user\"}[5m])) * 100"
      }
    }

    recipients {
      emails                    = [var.notification_email]
      notification_endpoint_ids = var.notification_endpoint_id != null ? [var.notification_endpoint_id] : null
    }
  }
}

# Case 2: Two Queries, No AI
# Math expression comparing two queries
resource "logzio_unified_alert" "metric_two_queries_no_ai" {
  title       = "Test: Two Queries No AI - A > B Comparison"
  type        = "METRIC_ALERT"
  description = "Alert when query A exceeds query B using math expression - no RCA"
  tags        = ["test", "metrics", "math", "comparison", "no-ai"]
  enabled     = true

  dashboard_id = var.dashboard_id != "" ? var.dashboard_id : null
  folder_id    = var.folder_id != "" ? var.folder_id : null
  panel_id = var.folder_id


  metric_alert {
    severity = "MEDIUM"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "$A - $B"
      metric_operator          = "ABOVE"
      min_threshold            = 0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "sum(rate(http_requests_total{status=\"200\"}[5m]))"
      }
    }

    queries {
      ref_id = "B"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "sum(rate(http_requests_total{status=\"500\"}[5m]))"
      }
    }

    recipients {
      emails                    = [var.notification_email]
      notification_endpoint_ids = var.notification_endpoint_id != null ? [var.notification_endpoint_id] : null
    }
  }
}

# Case 3: One Query, AI (RCA with same endpoints)
# Single query with RCA enabled, using alert notification endpoints for RCA
resource "logzio_unified_alert" "metric_one_query_with_ai" {
  title       = "Test: One Query With AI - High Memory"
  type        = "METRIC_ALERT"
  description = "Alert on high memory usage with RCA enabled - same notification endpoints"
  tags        = ["test", "metrics", "memory", "ai", "rca"]
  enabled     = true

  dashboard_id = var.dashboard_id != "" ? var.dashboard_id : null
  folder_id    = var.folder_id != "" ? var.folder_id : null
  panel_id = var.folder_id

  # RCA (AI) configuration - use same endpoints as alert
  rca                                      = true
  use_alert_notification_endpoints_for_rca = true
  runbook                                  = "1. Check memory usage trends\n2. Identify top consumers\n3. Review recent deployments\n4. Scale resources if needed"

  metric_alert {
    severity = "HIGH"

    trigger {
      trigger_type             = "THRESHOLD"
      metric_operator          = "ABOVE"
      min_threshold            = 85.0
      search_timeframe_minutes = 10
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "avg(node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100"
      }
    }

    recipients {
      emails                    = [var.notification_email]
      notification_endpoint_ids = var.notification_endpoint_id != null ? [var.notification_endpoint_id] : null
    }
  }
}

# Case 4: Two Queries, AI (RCA with same endpoints)
# Math expression with RCA enabled, using alert notification endpoints
resource "logzio_unified_alert" "metric_two_queries_with_ai" {
  title       = "Test: Two Queries With AI - Error Rate Percentage"
  type        = "METRIC_ALERT"
  description = "Alert when error rate exceeds threshold using math - RCA enabled with same endpoints"
  tags        = ["test", "metrics", "errors", "math", "ai", "rca"]
  enabled     = true

  dashboard_id = var.dashboard_id != "" ? var.dashboard_id : null
  folder_id    = var.folder_id != "" ? var.folder_id : null
  panel_id = var.folder_id

  # RCA configuration - use same endpoints as alert
  rca                                      = true
  use_alert_notification_endpoints_for_rca = true
  runbook                                  = "1. Check error logs for details\n2. Review affected endpoints\n3. Compare with baseline\n4. Initiate incident response if sustained"

  metric_alert {
    severity = "SEVERE"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "($A / $B) * 100"
      metric_operator          = "ABOVE"
      min_threshold            = 5.0
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
      emails                    = [var.notification_email]
      notification_endpoint_ids = var.notification_endpoint_id != null ? [var.notification_endpoint_id] : null
    }
  }
}

# =============================================================================
# OPTIONAL METRIC ALERT TEST CASES (Separate RCA endpoints)
# =============================================================================

# Case 5: One Query, AI with different endpoints
# Single query with RCA using separate notification endpoints
resource "logzio_unified_alert" "metric_one_query_ai_separate_endpoints" {
  title       = "Test: One Query AI Separate Endpoints - Disk Usage"
  type        = "METRIC_ALERT"
  description = "Alert on high disk usage with RCA using separate notification endpoints"
  tags        = ["test", "metrics", "disk", "ai", "rca", "separate-endpoints"]
  enabled     = true

  dashboard_id = var.dashboard_id != "" ? var.dashboard_id : null
  folder_id    = var.folder_id != "" ? var.folder_id : null
  panel_id = var.folder_id

  # RCA configuration - use different endpoints for RCA
  rca                                      = true
  use_alert_notification_endpoints_for_rca = false
  rca_notification_endpoint_ids            = var.rca_notification_endpoint_id != null ? [var.rca_notification_endpoint_id] : null
  runbook                                  = "1. Check disk usage by mount point\n2. Identify large files/directories\n3. Clean up temporary files\n4. Expand storage if needed"

  metric_alert {
    severity = "MEDIUM"

    trigger {
      trigger_type             = "THRESHOLD"
      metric_operator          = "ABOVE"
      min_threshold            = 90.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = var.datasource_uid
        promql_query   = "avg(node_filesystem_avail_bytes{mountpoint=\"/\"} / node_filesystem_size_bytes{mountpoint=\"/\"}) * 100"
      }
    }

    recipients {
      emails                    = [var.notification_email]
      notification_endpoint_ids = var.notification_endpoint_id != null ? [var.notification_endpoint_id] : null
    }
  }
}

# Case 6: Two Queries, AI with different endpoints
# Math expression with RCA using separate notification endpoints
resource "logzio_unified_alert" "metric_two_queries_ai_separate_endpoints" {
  title       = "Test: Two Queries AI Separate Endpoints - Latency vs Baseline"
  type        = "METRIC_ALERT"
  description = "Alert when latency deviation exceeds baseline - RCA with separate endpoints"
  tags        = ["test", "metrics", "latency", "math", "ai", "rca", "separate-endpoints"]
  enabled     = true

  dashboard_id = var.dashboard_id != "" ? var.dashboard_id : null
  folder_id    = var.folder_id != "" ? var.folder_id : null
  panel_id = var.folder_id


  # RCA configuration - use different endpoints for RCA
  rca                                      = true
  use_alert_notification_endpoints_for_rca = false
  rca_notification_endpoint_ids            = var.rca_notification_endpoint_id != null ? [var.rca_notification_endpoint_id] : null
  runbook                                  = "1. Review latency trends\n2. Check for degraded services\n3. Compare with historical baselines\n4. Investigate dependencies"

  metric_alert {
    severity = "HIGH"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "(($A - $B) / $B) * 100"
      metric_operator          = "ABOVE"
      min_threshold            = 50.0
      search_timeframe_minutes = 10
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
      emails                    = [var.notification_email]
      notification_endpoint_ids = var.notification_endpoint_id != null ? [var.notification_endpoint_id] : null
    }
  }
}

# Data source to test reading an alert
data "logzio_unified_alert" "test_read" {
  type     = "LOG_ALERT"
  alert_id = logzio_unified_alert.test_log_alert.alert_id
}

# =============================================================================
# OUTPUTS
# =============================================================================

# Log Alert Outputs
output "log_alert_id" {
  value       = logzio_unified_alert.test_log_alert.alert_id
  description = "ID of the test log alert"
}

output "log_alert_created_at" {
  value       = logzio_unified_alert.test_log_alert.created_at
  description = "Creation timestamp"
}

output "data_source_title" {
  value       = data.logzio_unified_alert.test_read.title
  description = "Title from data source"
}

# Metric Alert Outputs - Prioritized Cases
output "metric_one_query_no_ai_id" {
  value       = logzio_unified_alert.metric_one_query_no_ai.alert_id
  description = "Case 1: One Query No AI - Alert ID"
}

output "metric_two_queries_no_ai_id" {
  value       = logzio_unified_alert.metric_two_queries_no_ai.alert_id
  description = "Case 2: Two Queries No AI - Alert ID"
}

output "metric_one_query_with_ai_id" {
  value       = logzio_unified_alert.metric_one_query_with_ai.alert_id
  description = "Case 3: One Query With AI - Alert ID"
}

output "metric_one_query_with_ai_rca_enabled" {
  value       = logzio_unified_alert.metric_one_query_with_ai.rca
  description = "Case 3: RCA enabled status"
}

output "metric_two_queries_with_ai_id" {
  value       = logzio_unified_alert.metric_two_queries_with_ai.alert_id
  description = "Case 4: Two Queries With AI - Alert ID"
}

output "metric_two_queries_with_ai_rca_enabled" {
  value       = logzio_unified_alert.metric_two_queries_with_ai.rca
  description = "Case 4: RCA enabled status"
}

# Metric Alert Outputs - Optional Cases (Separate Endpoints)
output "metric_one_query_ai_separate_endpoints_id" {
  value       = logzio_unified_alert.metric_one_query_ai_separate_endpoints.alert_id
  description = "Case 5: One Query AI Separate Endpoints - Alert ID"
}

output "metric_two_queries_ai_separate_endpoints_id" {
  value       = logzio_unified_alert.metric_two_queries_ai_separate_endpoints.alert_id
  description = "Case 6: Two Queries AI Separate Endpoints - Alert ID"
}

# Summary Output
output "test_summary" {
  value = {
    log_alerts_created    = 1
    metric_alerts_created = 6
    total_alerts_created  = 7
    prioritized_cases     = 4
    optional_cases        = 2
  }
  description = "Summary of all created test alerts"
}

