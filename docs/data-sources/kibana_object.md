# Kibana Object Datasource

Provides a Logz.io kibana object datasource.

* Learn more about kibana objects in the [Logz.io Docs](https://docs.logz.io/api/#tag/Import-or-export-Kibana-objects)

## Example Usage

```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_kibana_object" "my_kb_obj" {
  kibana_version = "7.2.1"
  data = file("/path/to/your/object/config.json")
}

data "logzio_kibana_object" "my_kb_obj_datasource" {
  object_id = "search:tf-provider-test-search"
  object_type = "search"
  depends_on = ["logzio_kibana_object.my_kb_obj"]
}
```

## Argument Reference

* `object_id` - (String) The id of the Kibana Object.
* `object_type` - (String) The type of the Kibana Object. Can be one of the following: `search`, `dashboard`, `visualization`.

## Attribute Reference

* `kibana_version` - (String) The version of Kibana used at the time of export.
* `data` - (String) Exported Kibana objects.
