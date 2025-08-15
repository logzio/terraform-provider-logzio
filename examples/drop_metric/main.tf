terraform {
    required_providers {
        logzio = {
            source  = "logzio/logzio"
            version = "~> 1.20"
        }
    }
}

# Provider reads creds from env: LOGZIO_API_TOKEN, LOGZIO_REGION
provider "logzio" {}

# Variables
variable "metrics_account_id" {
    description = "Target Logz.io Metrics account ID"
    type        = number
}

variable "operator" {
    description = "Logical operator to combine filters"
    type        = string
    default     = "AND"
}

variable "drop_filters" {
    description = <<EOT
List of drop filters to create.
Each filter is a map with:
    key         = friendly name (must be unique)
    active      = bool
    filters     = list of { name, value, condition }
EOT
    type = map(object({
        active  = bool
        filters = list(object({
            name      = string
            value     = string
            condition = string
        }))
    }))
}

# Resource: create one drop filter per map entry
resource "logzio_drop_metrics" "this" {
    for_each    = var.drop_filters

    account_id  = var.metrics_account_id
    active      = each.value.active
    operator    = var.operator

    dynamic "filters" {
        for_each = each.value.filters
        content {
            name      = filters.value.name
            value     = filters.value.value
            condition = filters.value.condition
        }
    }
}

# Outputs
output "drop_filter_ids" {
    value = { for k, r in logzio_drop_metrics.this : k => r.id }
}
