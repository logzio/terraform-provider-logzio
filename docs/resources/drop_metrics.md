# Metrics Drop Filter Provider

Provides a Logz.io metrics drop filter resource. This can be used to create and manage Logz.io metrics drop filters.

* Learn more about drop filters in the [Logz.io Docs](https://docs.logz.io/docs/user-guide/data-hub/drop-filters/drop-fiters-metrics/).

## Example Usage
```hcl
variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_drop_metrics" "test_filter" {
  account_id = 1234
  filters {
    name = "__name__"
    value = "my_metric_name"
    condition = "EQ"
  }
  filters  {
    name = "my_label_name"
    value = "value_to_match"
    condition = "EQ"
  }
}
```

## Argument Reference
* `account_id` - (Required) The Logz.io metrics account ID to which the drop filter applies.
* `active` - (Optional) If true, the drop filter is active and metrics that match the filter are dropped before indexing. If false, the drop filter is disabled. Defaults to true.
* `filter` - (Required) The filter object that defines the drop filter criteria. See below for nested schema.
  * `name` - (Required) The name of the metric label to filter on.
  * `value` - (Required) The value of the metric label to match against.
  * `condition` - (Required) The comparison operator to use for the filter. Supported values are `EQ`, `NOT_EQ`, `REGEX_MATCH`, and `REGEX_NO_MATCH`.
* `operator` - (Optional) The logical operator for combining filter expressions. Supported value is `AND`.

## Attribute Reference
* `drop_metric_id` - (String) The unique identifier of the drop filter in the Logz.io database.
* `created_at` - (String) The timestamp when the drop filter was created.
* `updated_at` - (String) The timestamp when the drop filter was last updated.
* `created_by` - (String) The user who created the drop filter.
* `modified_at` - (String) The timestamp when the drop filter was last modified.
* `modified_by` - (String) The user who last updated the drop filter.

### Import metrics drop filter as resource

You can import drop filters as follows:

```
terraform import logzio_drop_metrics.my_filter <DROP-FILTER-ID>
```
