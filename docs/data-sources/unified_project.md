# Unified Project Data Source

Use this data source to read an existing Logz.io unified project (Perses project/dashboard folder) by its stable project ID or its Perses metadata name.

## Example Usage

```hcl
data "logzio_unified_project" "by_id" {
  id = "project-id"
}

data "logzio_unified_project" "by_metadata_name" {
  name = "system-metrics"
}
```

## Argument Reference

Configure exactly one of:

* `id` - (Optional, String) The project's stable unique identifier (GUID).
* `name` - (Optional, String) The project's Perses `metadata.name`. This selector does not match the human-readable display name.

## Attribute Reference

* `id` - (String) The stable unique project identifier. This is also used as Terraform's internal data-source ID.
* `name` - (String) The project's Perses `metadata.name`.
* `display_name` - (String) The human-readable project display name.
* `description` - (String) The project description.
