# logzio_metrics_rollup_rules

Provides a Logz.io metrics rollup rules resource. This allows you to manage metrics rollup rules in your Logz.io account.

## Example Usage

```hcl
resource "logzio_metrics_rollup_rules" "cpu_usage_rollup" {
  account_id               = 123456
  metric_name              = "cpu_usage"
  metric_type              = "gauge"
  rollup_function          = "last"
  labels_elimination_method = "exclude_by"
  labels                   = ["instance_id", "process_id"]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) The metrics account ID for the metrics rollup rule.
* `metric_name` - (Required) The name of the metric for which to create the rollup rule.
* `metric_type` - (Required) The type of the metric. Valid values are `gauge` and `counter`.
* `rollup_function` - (Required) The aggregation function to use for rolling up the metric. Valid values are `sum`, `min`, `max`, `count`, and `last`.
* `labels_elimination_method` - (Required) The method for eliminating labels. Valid values are `exclude_by`.
* `labels` - (Required) A list of label names to be eliminated from the metric.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the metrics rollup rule.

## Import

Metrics rollup rules can be imported using their ID:

```bash
terraform import logzio_metrics_rollup_rules.my_rollup_rule "rule_id"
``` 