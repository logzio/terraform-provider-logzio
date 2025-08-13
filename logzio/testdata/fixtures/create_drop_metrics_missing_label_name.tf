resource "logzio_drop_metrics" "%s" {
  account_id = %s
  filters {
    value = "my_metric"
    condition = "EQ"
  }
}
