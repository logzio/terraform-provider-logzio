# logzio_metrics_rollup_rules

Use this data source to access information about an existing Logz.io metrics rollup rule.

## Example Usage

```hcl
# Get metrics rollup rule by ID
data "logzio_metrics_rollup_rules" "my_rollup_rule" {
  id = "rule_id"
}

# Data source requires rule ID
# data "logzio_metrics_rollup_rules" "cpu_rollup" {
#   id = "rule_id_from_resource"
# }
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the metrics rollup rule.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the metrics rollup rule.
* `account_id` - The account ID of the metrics rollup rule.
* `metric_name` - The name of the metric.
* `metric_type` - The type of the metric.
* `rollup_function` - The aggregation function used for rolling up the metric.
* `labels_elimination_method` - The method for eliminating labels.
* `labels` - The list of label names being eliminated from the metric. 