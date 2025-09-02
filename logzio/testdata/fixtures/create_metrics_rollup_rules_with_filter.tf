resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_type = "COUNTER"
  labels_elimination_method = "GROUP_BY"
  labels = ["service", "region"]
  
  filter {
    expression {
      comparison = "EQ"
      name = "service"
      value = "frontend"
    }
    expression {
      comparison = "REGEX_MATCH"
      name = "region"
      value = "us-.*"
    }
  }
  
  name = "frontend_metrics_rollup"
  new_metric_name_template = "rollup.frontend.${metric_name}"
  drop_original_metric = true
} 