terraform {
  required_providers {
    logzio = {
      source = "logzio/logzio"
    }
  }
}

provider "logzio" {
  api_token = var.api_token
  region    = var.region
}

variable "api_token" {
  type        = string
  description = "Logz.io API token"
}

variable "region" {
  type        = string
  description = "Logz.io region (e.g., us, eu)"
  default     = "us"
}

variable "account_id" {
  type        = number
  description = "Account ID for the metrics rollup rule"
}

# Basic metric name-based rollup rule
resource "logzio_metrics_rollup_rules" "cpu_usage_rollup" {
  account_id                = var.account_id
  name                      = "CPU Usage Rollup"
  metric_name               = "cpu_usage"
  metric_type               = "GAUGE"
  rollup_function           = "LAST"
  labels_elimination_method = "EXCLUDE_BY"
  labels                    = ["instance_id", "process_id"]
}

# Counter with aggregation
resource "logzio_metrics_rollup_rules" "request_count_rollup" {
  account_id                = var.account_id
  name                      = "HTTP Requests Rollup"
  metric_name               = "http_requests_total"
  metric_type               = "COUNTER"
  rollup_function           = "SUM"
  labels_elimination_method = "EXCLUDE_BY"
  labels                    = ["path", "user_agent"]
}

# Filter-based rule with advanced features
resource "logzio_metrics_rollup_rules" "frontend_service_rollup" {
  account_id                = var.account_id
  name                      = "Frontend Service Metrics"
  metric_type               = "COUNTER"
  rollup_function           = "SUM"
  labels_elimination_method = "GROUP_BY"
  labels                    = ["service", "region", "environment"]
  
  filter {
    expression {
      comparison = "EQ"
      name       = "service"
      value      = "frontend"
    }
    expression {
      comparison = "REGEX_MATCH"
      name       = "environment"
      value      = "(prod|staging)"
    }
  }
  
  new_metric_name_template = "rollup.frontend.$${metric_name}"
  drop_original_metric     = true
}

# MEASUREMENT type for response time percentiles
resource "logzio_metrics_rollup_rules" "response_time_p95" {
  account_id                = var.account_id
  name                      = "Response Time P95"
  metric_name               = "http_response_time"
  metric_type               = "MEASUREMENT"
  rollup_function           = "P95"
  labels_elimination_method = "EXCLUDE_BY"
  labels                    = ["endpoint", "method", "status_code"]
}

# DELTA_COUNTER example
resource "logzio_metrics_rollup_rules" "error_rate_rollup" {
  account_id                = var.account_id
  name                      = "Error Rate Rollup"
  metric_name               = "errors_delta"
  metric_type               = "DELTA_COUNTER"
  rollup_function           = "SUM"
  labels_elimination_method = "GROUP_BY"
  labels                    = ["service", "region"]
}

# Data source examples
data "logzio_metrics_rollup_rules" "existing_cpu_rollup" {
  id = logzio_metrics_rollup_rules.cpu_usage_rollup.id
}

data "logzio_metrics_rollup_rules" "existing_frontend_rollup" {
  id = logzio_metrics_rollup_rules.frontend_service_rollup.id
}

# Outputs showcasing new attributes
output "cpu_rollup_rule_id" {
  value = logzio_metrics_rollup_rules.cpu_usage_rollup.id
}

output "cpu_rollup_details" {
  value = {
    id                        = data.logzio_metrics_rollup_rules.existing_cpu_rollup.id
    name                      = data.logzio_metrics_rollup_rules.existing_cpu_rollup.name
    metric_name               = data.logzio_metrics_rollup_rules.existing_cpu_rollup.metric_name
    metric_type               = data.logzio_metrics_rollup_rules.existing_cpu_rollup.metric_type
    rollup_function           = data.logzio_metrics_rollup_rules.existing_cpu_rollup.rollup_function
    labels_elimination_method = data.logzio_metrics_rollup_rules.existing_cpu_rollup.labels_elimination_method
    labels                    = data.logzio_metrics_rollup_rules.existing_cpu_rollup.labels
    version                   = data.logzio_metrics_rollup_rules.existing_cpu_rollup.version
  }
}

output "frontend_rollup_details" {
  value = {
    id                         = data.logzio_metrics_rollup_rules.existing_frontend_rollup.id
    name                       = data.logzio_metrics_rollup_rules.existing_frontend_rollup.name
    metric_type                = data.logzio_metrics_rollup_rules.existing_frontend_rollup.metric_type
    rollup_function            = data.logzio_metrics_rollup_rules.existing_frontend_rollup.rollup_function
    labels_elimination_method  = data.logzio_metrics_rollup_rules.existing_frontend_rollup.labels_elimination_method
    labels                     = data.logzio_metrics_rollup_rules.existing_frontend_rollup.labels
    new_metric_name_template   = data.logzio_metrics_rollup_rules.existing_frontend_rollup.new_metric_name_template
    drop_original_metric       = data.logzio_metrics_rollup_rules.existing_frontend_rollup.drop_original_metric
    filter                     = data.logzio_metrics_rollup_rules.existing_frontend_rollup.filter
    version                    = data.logzio_metrics_rollup_rules.existing_frontend_rollup.version
  }
}

output "all_rollup_rules" {
  description = "Summary of all created rollup rules"
  value = {
    cpu_usage = {
      id   = logzio_metrics_rollup_rules.cpu_usage_rollup.id
      name = logzio_metrics_rollup_rules.cpu_usage_rollup.name
      type = "metric_name_based"
    }
    request_count = {
      id   = logzio_metrics_rollup_rules.request_count_rollup.id
      name = logzio_metrics_rollup_rules.request_count_rollup.name
      type = "metric_name_based"
    }
    frontend_service = {
      id   = logzio_metrics_rollup_rules.frontend_service_rollup.id
      name = logzio_metrics_rollup_rules.frontend_service_rollup.name
      type = "filter_based"
    }
    response_time_p95 = {
      id   = logzio_metrics_rollup_rules.response_time_p95.id
      name = logzio_metrics_rollup_rules.response_time_p95.name
      type = "measurement_percentile"
    }
    error_rate = {
      id   = logzio_metrics_rollup_rules.error_rate_rollup.id
      name = logzio_metrics_rollup_rules.error_rate_rollup.name
      type = "delta_counter"
    }
  }
} 