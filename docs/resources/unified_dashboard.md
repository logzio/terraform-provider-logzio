# Unified Dashboard Provider

Provides a Logz.io unified dashboard resource. This can be used to create and manage Perses-based unified dashboards in Logz.io.

* Learn more about Logz.io's unified dashboards API in [Logz.io Docs](https://api-docs.logz.io/docs/logz/create-a-new-dashboard/).

## Example Usage

```hcl
resource "logzio_unified_dashboard" "my_dashboard" {
  folder_id = "my-project-id"
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

### Required:

* `folder_id` - (String) The unique identifier of the unified dashboard folder (Perses project) that stores the dashboard. Changing this value forces a new resource.
* `dashboard_json` - (String) The complete Perses dashboard document in JSON format. The document must include `kind`, `metadata.name`, and `spec`. Once created, you cannot change `metadata.name`.

## Attribute Reference

* `dashboard_uid` - (String) The stable unique identifier of the dashboard. Persist this value; do not use the API's version-row `id`.
* `name` - (String) The dashboard's Perses identity (`metadata.name`).
* `version` - (Int) Dashboard version.

### Import Logz.io Unified Dashboard as Terraform resource

You can import an existing dashboard as follows:

```
terraform import logzio_unified_dashboard.my_dashboard <FOLDER-ID>/<DASHBOARD-UID>
```
