metrics_account_id = 123456

drop_filters = {
  # Example 1: drop specified calls_total metric
  calls_total_filter = {
    active = true
    filters = [
      { name = "__name__", value = "calls_total", condition = "EQ" }
    ]
  }

  # Example 2: drop all latency span metrics in prod namespace
  latency_bucket_prod = {
    active = true
    filters = [
      { name = "__name__", value = "^latency.*", condition = "REGEX_MATCH" }
    ]
  }
}
