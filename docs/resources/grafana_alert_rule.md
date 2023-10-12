# Grafana Alert Rule Resource

Provides a Logz.io Grafana alert rule resource. This can be used to create and manage Grafana alert rules in Logz.io.

* Learn more about Logz.io's Grafana alert rule API in [Logz.io Docs](https://docs.logz.io/api/#tag/Grafana-alerting-provisioning).

## Example Usage

```hcl
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
  folder_uid = "my-folder-uid"
  for = 3
  no_data_state = "OK"
  rule_group = "rule_group_1"
  title = "my_grafana_alert"
}
```

## Argument Reference

### Required:

* `condition` - (String) The `ref_id` of the query node in the `data` field to use as the alert condition.
* `data` - (Block List) A sequence of stages that describe the contents of the rule. See below for **nested schema**.
