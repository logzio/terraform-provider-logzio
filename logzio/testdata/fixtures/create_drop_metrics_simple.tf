resource "logzio_drop_metrics" "%s" {
  account_id = %s
  filters {
    name = "__name__"
    value = "my_metric"
    condition = "EQ"
  }
}
