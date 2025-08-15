metrics_account_id = 123456

drop_filters = {
  # Example 1: drop calls_total
  calls_total_filter = {
    active = true
    filters = [
      { name = "__name__", value = "calls_total", condition = "EQ" }
    ]
  }

  # Example 2: drop latency_bucket in prod namespace
  latency_bucket_prod = {
    active = true
    filters = [
      { name = "__name__", value = ".*latency.*", condition = "REGEX_MATCH" }
    ]
  }
}
