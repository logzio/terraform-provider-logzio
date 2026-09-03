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

provider "logzio" {
  api_token = var.api_token
}

# A unified project acts as the folder that contains unified dashboards.
resource "logzio_unified_project" "system_metrics" {
  name        = "system-metrics"
  description = "Dashboards for system health and performance"
}

resource "logzio_unified_dashboard" "cpu_usage" {
  folder_id      = logzio_unified_project.system_metrics.folder_id
  dashboard_json = <<EOD
{
  "kind": "Dashboard",
  "metadata": {
    "name": "cpu-usage"
  },
  "spec": {
    "display": {
      "name": "CPU Usage"
    },
    "duration": "1h",
    "panels": {},
    "layouts": []
  }
}
EOD
}

# Read the created project and dashboard through their data sources.
data "logzio_unified_project" "system_metrics" {
  id = logzio_unified_project.system_metrics.folder_id
}

data "logzio_unified_dashboard" "cpu_usage" {
  folder_id     = data.logzio_unified_project.system_metrics.id
  dashboard_uid = logzio_unified_dashboard.cpu_usage.dashboard_uid
}

output "project" {
  value = {
    id           = data.logzio_unified_project.system_metrics.id
    name         = data.logzio_unified_project.system_metrics.name
    display_name = data.logzio_unified_project.system_metrics.display_name
  }
}

output "dashboard" {
  value = {
    uid  = data.logzio_unified_dashboard.cpu_usage.dashboard_uid
    name = data.logzio_unified_dashboard.cpu_usage.name
    json = data.logzio_unified_dashboard.cpu_usage.dashboard_json
  }
}
