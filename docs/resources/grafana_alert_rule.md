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
  folder_uid = "my-folder-uid"
  for = "3m"
  no_data_state = "OK"
  rule_group = "rule_group_1"
  title = "my_grafana_alert"
}
```

## Argument Reference

### Required:

* `condition` - (String) The `ref_id` of the query node in the `data` field to use as the alert condition.
* `data` - (Block List) A sequence of stages that describe the contents of the rule. See below for **nested schema**.
* `folder_uid` - (String) The UID of the folder that the alert rule belongs to.
* `for` - (String) The amount of time for which the rule must be breached for the rule to be considered to be Firing. Before this time has elapsed, the rule is only considered to be Pending. Should be in a duration string format, for example "3m0s".
* `rule_group` - (String) The rule group this rule is associated with.
* `title` - (String) The title of this rule.

### Optional: 

* `annotations` - (Map of String) Key-value pairs of metadata to attach to the alert rule that may add user-defined context, but cannot be used for matching, grouping, or routing.
* `labels` - (Map of String) Key-value pairs to attach to the alert rule that can be used in matching, grouping, and routing.
* `is_paused` - (Boolean) Sets whether the alert should be paused or not. Defaults to `false`.
* `no_data_state` - (String) Describes what state to enter when the rule's query returns No Data. Options are `OK`, `NoData`, and `Alerting`. Defaults to `NoData`.

#### Nested schema for `data`:

##### Required:

* `ref_id` - (String) A unique string to identify this query stage within a rule.
* `datasource_uid` - (String) The UID of the datasource being queried, or "-100" if this stage is an expression stage.
* `model` - (String) Custom JSON data to send to the specified datasource when querying.
* `relative_time_range` - (Block List, Min: 1, Max: 1) The time range, relative to when the query is executed, across which to query. See below for **nested schema**.

##### Optional:

* `query_type` - (String) An optional identifier for the type of query being executed.

#### Nested schema for `data.relative_time_range`:

##### Required:

* `from` - (Integer) The number of seconds in the past, relative to when the rule is evaluated, at which the time range begins.
* `to` - (Integer) The number of seconds in the past, relative to when the rule is evaluated, at which the time range ends.
