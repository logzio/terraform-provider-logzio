resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "http_requests_total"
  metric_type = "COUNTER"
  rollup_function = "LAST"
  labels_elimination_method = "EXCLUDE_BY"
  labels = ["path"]
} 