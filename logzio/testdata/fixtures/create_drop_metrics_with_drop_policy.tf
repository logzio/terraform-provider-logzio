resource "logzio_drop_metrics" "%s" {
  account_id = %s
  name = "test-drop-metrics-with-policy"
  drop_policy = "DROP_BEFORE_STORING"
  
  filters {
    name = "__name__"
    value = "my_metric"
    condition = "EQ"
  }
} 