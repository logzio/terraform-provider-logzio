# Unified Project Resource

Provides a Logz.io unified project resource. A unified project is a Perses project that acts as a folder for unified dashboards.

## Example Usage

```hcl
resource "logzio_unified_project" "system_metrics" {
  name        = "system-metrics"
  description = "Dashboards for system health and performance"
}

resource "logzio_unified_dashboard" "cpu_usage" {
  folder_id = logzio_unified_project.system_metrics.folder_id
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
```

## Argument Reference

### Required

* `name` - (String) The project's Perses identity (`metadata.name`). Changing this value forces a new resource.

### Optional

* `display_name` - (String) An optional human-readable project name. When omitted, it defaults to `name`.
* `description` - (String) The human-readable project description. When omitted, the value returned by the API is stored in state.

## Attribute Reference

* `folder_id` - (String) The stable unique project identifier. Use this value as `folder_id` when creating or reading unified dashboards.

The Terraform resource ID is the same stable project identifier exposed as `folder_id`.

## Import

Import an existing unified project using its project ID:

```shell
terraform import logzio_unified_project.system_metrics <PROJECT-ID>
```
