resource "logzio_drop_metrics" "%s" {
  account_id = %s

  operator = "OR"

  filters {
    name = "__name__"
    value = "my_metric"
    condition = "EQ"
  }

  filters {
      name = "my_label"
      value = "my_value"
      condition = "EQ"
    }
}
