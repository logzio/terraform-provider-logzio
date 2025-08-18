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

resource "logzio_metrics_rollup_rules" "cpu_usage_rollup" {
  account_id               = var.account_id
  metric_name              = "cpu_usage"
  metric_type              = "GAUGE"
  rollup_function          = "LAST"
  labels_elimination_method = "EXCLUDE_BY"
  labels                   = ["instance_id", "process_id"]
}

resource "logzio_metrics_rollup_rules" "request_count_rollup" {
  account_id               = var.account_id
  metric_name              = "http_requests_total"
  metric_type              = "COUNTER"
  rollup_function          = "SUM"
  labels_elimination_method = "EXCLUDE_BY"
  labels                   = ["path", "user_agent"]
}

# Data source example
data "logzio_metrics_rollup_rules" "existing_cpu_rollup" {
  id = logzio_metrics_rollup_rules.cpu_usage_rollup.id
}

output "cpu_rollup_rule_id" {
  value = logzio_metrics_rollup_rules.cpu_usage_rollup.id
}

output "existing_cpu_rollup" {
  value = data.logzio_metrics_rollup_rules.existing_cpu_rollup
} 