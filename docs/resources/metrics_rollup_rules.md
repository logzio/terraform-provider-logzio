# logzio_metrics_rollup_rules

Provides a Logz.io metrics rollup rules resource. This allows you to manage metrics rollup rules in your Logz.io account.

## Example Usage

### Basic metric name-based rule

```hcl
resource "logzio_metrics_rollup_rules" "cpu_usage_rollup" {
  account_id                = 123456
  metric_name               = "cpu_usage"
  metric_type               = "GAUGE"
  rollup_function           = "LAST"
  labels_elimination_method = "EXCLUDE_BY"
  labels                    = ["instance_id", "process_id"]
}
```

### Filter-based rule with advanced features

```hcl
resource "logzio_metrics_rollup_rules" "frontend_metrics_rollup" {
  account_id                = 123456
  name                      = "Frontend Service Metrics"
  metric_type               = "COUNTER"
  labels_elimination_method = "GROUP_BY"
  labels                    = ["service", "region"]
  
  filter {
    expression {
      comparison = "EQ"
      name       = "service"
      value      = "frontend"
    }
    expression {
      comparison = "REGEX_MATCH"
      name       = "region"
      value      = "us-.*"
    }
  }
  
  new_metric_name_template = "rollup.frontend.${metric_name}"
  drop_original_metric     = true
}
```

### MEASUREMENT metric type with statistical aggregation

```hcl
resource "logzio_metrics_rollup_rules" "response_time_rollup" {
  account_id                = 123456
  metric_name               = "http_response_time"
  metric_type               = "MEASUREMENT"
  rollup_function           = "P95"
  labels_elimination_method = "EXCLUDE_BY"
  labels                    = ["endpoint", "method"]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) The metrics account ID for the metrics rollup rule.
* `metric_name` - (Optional) The name of the metric for which to create the rollup rule. Either `metric_name` or `filter` must be specified, but not both.
* `metric_type` - (Required) The type of the metric. Valid values are `GAUGE`, `COUNTER`, `DELTA_COUNTER`, `CUMULATIVE_COUNTER`, and `MEASUREMENT`.
* `rollup_function` - (Optional) The aggregation function to use for rolling up the metric. Required for `GAUGE` and `MEASUREMENT` metric types, not allowed for `COUNTER`, `DELTA_COUNTER`, and `CUMULATIVE_COUNTER` types. Valid values include `SUM`, `MIN`, `MAX`, `COUNT`, `LAST`, `MEAN`, `MEDIAN`, `STDEV`, `SUMSQ`, and percentiles (`P10`, `P20`, `P25`, `P30`, `P40`, `P50`, `P60`, `P70`, `P75`, `P80`, `P90`, `P95`, `P99`, `P999`, `P9999`). Note: For `MEASUREMENT` metric type, only `SUM`, `MIN`, `MAX`, `COUNT`, `SUMSQ`, `MEAN`, and `LAST` are allowed.
* `labels_elimination_method` - (Required) The method for eliminating labels. Valid values are `EXCLUDE_BY` and `GROUP_BY`.
* `labels` - (Required) A list of label names to be eliminated from the metric.
* `name` - (Optional) A human-readable name for the rollup rule.
* `filter` - (Optional) A filter block to match metrics by label values. Either `metric_name` or `filter` must be specified, but not both.
  * `expression` - (Required) A list of filter expressions.
    * `comparison` - (Required) The comparison operator. Valid values are `EQ`, `NOT_EQ`, `REGEX_MATCH`, and `REGEX_NO_MATCH`.
    * `name` - (Required) The label name to match against.
    * `value` - (Required) The value to match.
* `new_metric_name_template` - (Optional) A template for generating new metric names. Use `${metric_name}` to reference the original metric name.
* `drop_original_metric` - (Optional) Whether to drop the original metric after creating the rollup. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the metrics rollup rule.

## Import

Metrics rollup rules can be imported using their ID:

```bash
terraform import logzio_metrics_rollup_rules.my_rollup_rule "rule_id"
``` 