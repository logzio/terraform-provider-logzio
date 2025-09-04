resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "response_time"
  metric_type = "MEASUREMENT"
  rollup_function = "MEAN"
  labels_elimination_method = "EXCLUDE_BY"
  labels = ["instance_id", "endpoint"]
} 