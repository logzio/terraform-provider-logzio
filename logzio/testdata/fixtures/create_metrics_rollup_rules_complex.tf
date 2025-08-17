resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "http_requests_total"
  metric_type = "counter"
  rollup_function = "sum"
  labels_elimination_method = "exclude_by"
  labels = ["path", "method", "user_agent"]
} 