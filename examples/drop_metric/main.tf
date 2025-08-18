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

variable "account_id" {
  type        = number
  description = "Account ID for the drop metrics rules"
}

provider "logzio" {
  api_token = var.api_token
}

# Simple drop metrics example - drop specific metric by name
resource "logzio_drop_metrics" "simple_metric_drop" {
  account_id = var.account_id
  
  filters {
    name      = "__name__"
    value     = "cpu_usage_percent"
    condition = "EQ"
  }
}

# Complex drop metrics example - drop metrics based on multiple conditions
resource "logzio_drop_metrics" "complex_metric_drop" {
  account_id = var.account_id
  active     = true
  operator   = "AND"
  
  filters {
    name      = "__name__"
    value     = "http_requests_total"
    condition = "EQ"
  }
  
  filters {
    name      = "environment"
    value     = "staging"
    condition = "EQ"
  }
  
  filters {
    name      = "status_code"
    value     = "2xx"
    condition = "REGEX_MATCH"
  }
}

# Drop metrics using regex pattern matching
resource "logzio_drop_metrics" "regex_metric_drop" {
  account_id = var.account_id
  
  filters {
    name      = "__name__"
    value     = "test_.*"
    condition = "REGEX_MATCH"
  }
}

# Drop metrics using NOT_EQ condition
resource "logzio_drop_metrics" "exclude_metric_drop" {
  account_id = var.account_id
  
  filters {
    name      = "service"
    value     = "critical-service"
    condition = "NOT_EQ"
  }
}

# Data source example - retrieve existing drop metrics rule
data "logzio_drop_metrics" "existing_drop_rule" {
  drop_metric_id = logzio_drop_metrics.simple_metric_drop.drop_metric_id
}

# Outputs
output "simple_drop_rule_id" {
  description = "ID of the simple drop metrics rule"
  value       = logzio_drop_metrics.simple_metric_drop.drop_metric_id
}

output "complex_drop_rule_id" {
  description = "ID of the complex drop metrics rule"
  value       = logzio_drop_metrics.complex_metric_drop.drop_metric_id
}

output "existing_drop_rule_details" {
  description = "Details of the existing drop metrics rule"
  value       = data.logzio_drop_metrics.existing_drop_rule
} 