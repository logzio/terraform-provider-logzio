resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_type = "GAUGE"
  rollup_function = "SUM"
  labels_elimination_method = "EXCLUDE_BY"
  labels = ["label1"]
} 