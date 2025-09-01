resource "logzio_drop_metrics" "%s" {
  account_id = %s
  name = "test-drop-metrics-filter"
  
  filters {
    name = "__name__"
    value = "my_metric"
    condition = "EQ"
  }
} 