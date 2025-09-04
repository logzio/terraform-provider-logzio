resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "test_metric"
  metric_type = "GAUGE"
  rollup_function = "AVERAGE"
  labels_elimination_method = "EXCLUDE_BY"
  labels = ["label1"]
} 