# Drop Filter Provider

Provides a Logz.io drop filter resource. This can be used to create and manage Logz.io logs drop filters.

* Learn more about drop filters in the [Logz.io Docs](https://docs.logz.io/docs/user-guide/data-hub/drop-filters/drop-fiters-logs/).

## Example Usage

```hcl
variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_drop_filter" "test_filter" {
  log_type = "some_type"

  field_conditions {
    field_name = "some_field"
    value = "string_value"
  }
  field_conditions {
    field_name = "other_field_int"
    value = 200
  }
}
```

## Argument Reference

**Note:** Except the `active` argument, changing an argument value _after creation_ will cause the **entire resource** to be destroyed and re-created.

### Required:

* `field_conditions` - (Block list) Filters for an exact match of a field:value pair. **Note** that changing this field after creation will cause the resource to be destroyed and re-created. See below for **nested schema**.

### Optional:

* `log_type` - (String) Filters for the [log type](https://docs.logz.io/user-guide/log-shipping/built-in-log-types.html). Omit or leave empty if you want this filter to apply to all types. **Note** that changing this field after creation will cause the resource to be destroyed and re-created. 
* `active` - (Boolean) If true, the drop filter is active and logs that match the filter are dropped before indexing. If false, the drop filter is disabled. **Note** this argument can only be changed after the creation of the filter. Each filter is created with the `active` argument set to true.
* `gb_threshold` - (Float) The threshold in GB for the drop filter. If the total size of the logs that match the filter exceeds this threshold, the logs are dropped before indexing. If not specified, the default is `0`, which means that all logs that match the filter are dropped.

#### Nested schema for `field_conditions`:

* `field_name` - (String) Exact field name in your Kibana mapping for the selected `log_type`. **Note** that changing this field after creation will cause the resource to be destroyed and re-created.
* `value` - (Object) Exact field value. The filter looks for an exact value match of the entire object. **Note** that changing this field after creation will cause the resource to be destroyed and re-created.

##  Attribute Reference

* `drop_filter_id` - (String) Drop filter ID in the Logz.io database.

### Import drop filter as resource

You can import drop filters as follows:

```
terraform import logzio_drop_filter.my_filter <DROP-FILTER-ID>
```
