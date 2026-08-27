# Unified Dashboard Data Source

Use this data source to read an existing Perses-based unified dashboard from a specific Logz.io unified project.

## Example Usage

```hcl
data "logzio_unified_dashboard" "cpu_usage" {
  folder_id     = "my-project-id"
  dashboard_uid = "my-dashboard-uid"
}

output "dashboard_document" {
  value = data.logzio_unified_dashboard.cpu_usage.dashboard_json
}
```

## Argument Reference

* `folder_id` - (Required, String) The unique identifier of the unified project containing the dashboard.
* `dashboard_uid` - (Required, String) The dashboard's stable route identifier. This is distinct from the API's mutable version-row `id`.

The resource computes `dashboard_uid` after creating a dashboard. The data source requires it because both `folder_id` and `dashboard_uid` are needed to locate an existing dashboard.

The Terraform data-source ID is the stable composite `folder_id/dashboard_uid`.

## Attribute Reference

* `dashboard_json` - (String) The normalized Perses dashboard document. It uses the same normalization as the `logzio_unified_dashboard` resource: the document is preserved, but server-owned metadata is removed and only `metadata.name` is retained in `metadata`.
* `name` - (String) The dashboard's Perses `metadata.name`.
* `version` - (Int) The dashboard version.
* `folder_id` - (String) The requested unified project identifier.
* `dashboard_uid` - (String) The requested stable dashboard identifier.
