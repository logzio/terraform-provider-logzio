# Grafana Folder Provider

Provides a Logz.io Grafana folder resource. This can be used to create and manage Grafana folders in Logz.io.

* Learn more about Logz.io's Grafana folder API in [Logz.io Docs]().

## Example Usage

```hcl
resource logzio_grafana_folder "my_folder" {
  title = "another_title"
}
```

## Argument Reference

### Required:

- `title` - (String) The title of the folder.

## Attribute Reference

- `uid` - (String) Unique identifier for the folder.
- `folder_id` - (Integer) Auto-incrementing numeric value.
- `url` - (String) Url for the folder.
- `version` - (Integer) Version number of the folder.

### Import Logz.io Grafana folder as Terraform resource

You can import existing folder as follows:

```
terraform import logzio_grafana_folder.my_folder <FOLDER-ID>
```