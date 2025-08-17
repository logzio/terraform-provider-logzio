resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "test_metric"
  metric_type = "gauge"
  rollup_function = "sum"
  labels_elimination_method = "include_by"
  labels = ["label1"]
} 