resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "cpu_usage"
  metric_type = "gauge"
  rollup_function = "max"
  labels_elimination_method = "exclude_by"
  labels = ["instance_id", "region"]
} 