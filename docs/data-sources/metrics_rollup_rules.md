# logzio_metrics_rollup_rules

Use this data source to access information about an existing Logz.io metrics rollup rule.

## Example Usage

### Get by ID

```hcl
# Get metrics rollup rule by ID
data "logzio_metrics_rollup_rules" "my_rollup_rule" {
  id = "rule_id"
}

# Reference a rule created by resource
data "logzio_metrics_rollup_rules" "cpu_rollup" {
  id = logzio_metrics_rollup_rules.cpu_usage_rollup.id
}
```

### Search by attributes

```hcl
# Find rollup rule by account ID and metric name
data "logzio_metrics_rollup_rules" "cpu_rollup" {
  account_id  = 123456
  metric_name = "cpu_usage"
}

# Find rollup rule by name
data "logzio_metrics_rollup_rules" "frontend_rollup" {
  name = "Frontend Service Metrics"
}

# Find rollup rule by multiple criteria
data "logzio_metrics_rollup_rules" "specific_rollup" {
  account_id  = 123456
  metric_type = "GAUGE"
  name        = "CPU Usage Rollup"
}
```

### Output rule details

```hcl
output "rollup_rule_details" {
  value = {
    name                      = data.logzio_metrics_rollup_rules.my_rollup_rule.name
    metric_name               = data.logzio_metrics_rollup_rules.my_rollup_rule.metric_name
    metric_type               = data.logzio_metrics_rollup_rules.my_rollup_rule.metric_type
    rollup_function           = data.logzio_metrics_rollup_rules.my_rollup_rule.rollup_function
    labels_elimination_method = data.logzio_metrics_rollup_rules.my_rollup_rule.labels_elimination_method
    labels                    = data.logzio_metrics_rollup_rules.my_rollup_rule.labels
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the metrics rollup rule. If not specified, the data source will search for rules matching the other provided criteria.
* `account_id` - (Optional) Filter by the metrics account ID.
* `name` - (Optional) Filter by the rule name.
* `metric_name` - (Optional) Filter by the metric name.
* `metric_type` - (Optional) Filter by the metric type (GAUGE, COUNTER, DELTA_COUNTER, CUMULATIVE_COUNTER, or MEASUREMENT).

**Note**: If multiple rules match the search criteria, you must either specify an `id` or add more search criteria to uniquely identify a single rule.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the metrics rollup rule.
* `account_id` - The account ID of the metrics rollup rule.
* `name` - The human-readable name of the rollup rule.
* `metric_name` - The name of the metric (if rule is metric name-based).
* `metric_type` - The type of the metric (GAUGE, COUNTER, DELTA_COUNTER, CUMULATIVE_COUNTER, or MEASUREMENT).
* `rollup_function` - The aggregation function used for rolling up the metric. Always "SUM" for COUNTER, DELTA_COUNTER, and CUMULATIVE_COUNTER types.
* `labels_elimination_method` - The method for eliminating labels (EXCLUDE_BY or GROUP_BY).
* `labels` - The list of label names being eliminated from the metric. 
* `new_metric_name_template` - The template for generating new metric names.
* `drop_original_metric` - Whether the original metric is dropped after creating the rollup.
* `filter` - Filter configuration for rule matching (if rule is filter-based).
  * `expression` - List of filter expressions.
    * `comparison` - The comparison operator (EQ, NOT_EQ, REGEX_MATCH, or REGEX_NO_MATCH).
    * `name` - The label name being matched.
    * `value` - The value being matched.