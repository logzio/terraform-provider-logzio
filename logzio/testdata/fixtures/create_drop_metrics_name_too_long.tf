resource "logzio_drop_metrics" "%s" {
  account_id = %s
  name = "this-is-a-very-long-name-that-exceeds-256-characters-this-is-a-very-long-name-that-exceeds-256-characters-this-is-a-very-long-name-that-exceeds-256-characters-this-is-a-very-long-name-that-exceeds-256-characters-this-is-a-very-long-name-that-exceeds"
  
  filters {
    name = "__name__"
    value = "my_metric"
    condition = "EQ"
  }
} 