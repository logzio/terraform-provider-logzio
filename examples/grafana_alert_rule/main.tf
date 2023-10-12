resource "logzio_grafana_alert_rule" "test_grafana_alert" {
  annotations = {
    "foo" = "bar"
    "hello" = "world"
  }
  condition = "A"
  data {
    ref_id = "B"
    datasource_uid = "AB1C234567D89012E"
    query_type = ""
    model = jsonencode({
      hide          = false
      refId         = "B"
    })
    relative_time_range {
      from = 700
      to   = 0
    }
  }
  labels = {
    "hey" = "oh"
    "lets" = "go"
  }
  is_paused = false
  exec_err_state = "Alerting"
  folder_uid = "%s"
  for = 3
  no_data_state = "OK"
  rule_group = "rule_group_1"
  title = "my_grafana_alert"
}