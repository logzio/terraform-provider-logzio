resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "cpu_usage"
  metric_type = "GAUGE"
  rollup_function = "MAX"
  labels_elimination_method = "EXCLUDE_BY"
  labels = ["instance_id", "region"]
} 