# Unified Alert Resource

Provides a Logz.io unified alert resource. This resource allows you to create and manage both log-based and metric-based alerts through a single unified API.

The alert type is determined by the `type` field inside `alert_configuration`:
- When `type = "LOG_ALERT"`, configure log-alert fields (`sub_components`, `correlations`, `schedule`).
- When `type = "METRIC_ALERT"`, configure metric-alert fields (`trigger`, `queries`, `severity`).

## Example Usage - Log Alert

```hcl
resource "logzio_unified_alert" "log_alert_example" {
  title       = "High error rate in checkout service"
  description = "Triggers when the error rate of the checkout service exceeds the defined threshold."
  tags        = ["environment:production", "service:checkout"]
  enabled     = true

  linked_panel {
    folder_id    = "unified-folder-uid"
    dashboard_id = "unified-dashboard-uid"
    panel_id     = "A"
  }

  runbook = "If this alert fires, check checkout pods and logs, verify recent deployments, and roll back if necessary."

  rca                                      = true
  rca_notification_endpoint_ids            = [101, 102]
  use_alert_notification_endpoints_for_rca = true

  recipients {
    emails                    = ["devops@company.com", "oncall@company.com"]
    notification_endpoint_ids = [11, 12]
  }

  alert_configuration {
    type                           = "LOG_ALERT"
    search_timeframe_minutes       = 15
    suppress_notifications_minutes = 30
    alert_output_template_type     = "JSON"

    sub_components {
      query_definition {
        query = "kubernetes.container_name:checkout AND level:error"

        filters = jsonencode({
          bool = {
            must     = []
            should   = []
            filter   = []
            must_not = []
          }
        })

        group_by = ["kubernetes.pod_name"]

        aggregation {
          aggregation_type      = "SUM"
          field_to_aggregate_on = "error_count"
        }

        should_query_on_all_accounts = false
        account_ids_to_query_on      = [12345]
      }

      trigger {
        operator = "GREATER_THAN"

        severity_threshold_tiers {
          severity  = "HIGH"
          threshold = 100
        }

        severity_threshold_tiers {
          severity  = "MEDIUM"
          threshold = 50
        }
      }

      output {
        should_use_all_fields = false

        columns {
          field_name = "kubernetes.pod_name"
          sort       = "DESC"
        }
      }
    }

    correlations {
      correlation_operators = ["AND"]
    }

    schedule {
      cron_expression = "*/1 * * * *"
      timezone        = "UTC"
    }
  }
}
```

## Example Usage - Metric Alert (threshold)

```hcl
resource "logzio_unified_alert" "metric_alert_example" {
  title       = "High 5xx rate (absolute)"
  description = "Fire when 5xx requests exceed 5 req/min over 5 minutes."
  tags        = ["environment:production", "service:checkout"]
  enabled     = true

  linked_panel {
    folder_id    = "unified-folder-uid"
    dashboard_id = "unified-dashboard-uid"
    panel_id     = "A"
  }

  runbook = "RCA: inspect ingress errors by pod and compare to last deploy."

  recipients {
    emails                    = ["devops@company.com"]
    notification_endpoint_ids = [11, 12]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "INFO"

    trigger {
      type = "threshold"

      condition {
        operator_type = "above"
        threshold     = 5
      }
    }

    queries {
      ref_id = "A"

      query_definition {
        promql_query = "sum(rate(http_requests_total{status=~\"5..\"}[5m]))"
      }
    }
  }
}
```

## Example Usage - Metric Alert with Math Expression

```hcl
resource "logzio_unified_alert" "metric_math_alert" {
  title       = "5xx error rate percentage is high"
  description = "Fire when 5xx responses exceed 2% of total requests over 5 minutes."
  tags        = ["environment:production", "service:checkout"]
  enabled     = true

  recipients {
    emails                    = ["devops@company.com"]
    notification_endpoint_ids = [11, 12]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "INFO"

    trigger {
      type       = "math"
      expression = "(A / B) * 100"
    }

    queries {
      ref_id = "A"
      query_definition {
        promql_query = "sum(rate(http_requests_total{status=~\"5..\"}[5m]))"
      }
    }

    queries {
      ref_id = "B"
      query_definition {
        promql_query = "sum(rate(http_requests_total[5m]))"
      }
    }
  }
}
```

## Argument Reference

### Required Arguments

* `title` - (String) Alert name displayed in UI and notifications.
* `alert_configuration` - (Block, Max: 1) The alert configuration. See [Alert Configuration](#alert-configuration) below.

### Optional Arguments

* `description` - (String) Explanation of alert purpose and firing conditions.
* `tags` - (List of String) Labels for grouping and searching alerts.
* `linked_panel` - (Block, Max: 1) Dashboard panel link. See [Linked Panel](#linked-panel) below.
* `runbook` - (String) Operational instructions for responders; also used as RCA instruction text when RCA is enabled.
* `enabled` - (Boolean) Alert activation status. Default: `true`.
* `rca` - (Boolean) Enable Root Cause Analysis. Default: `false`.
* `rca_notification_endpoint_ids` - (List of Integer) Notification endpoint IDs for RCA results.
* `use_alert_notification_endpoints_for_rca` - (Boolean) When true, RCA uses same endpoints as alert. Default: `false`.
* `recipients` - (Block, Max: 1) Who receives notifications. See [Recipients](#recipients) below.

### Linked Panel

The `linked_panel` block supports:

* `folder_id` - (Optional, String) UID of the unified folder in Logz.io.
* `dashboard_id` - (Optional, String) UID of the unified dashboard for context linking.
* `panel_id` - (Optional, String) Specific panel ID on the dashboard.

### Recipients

The `recipients` block supports:

* `emails` - (Optional, List of String) Email addresses for notifications.
* `notification_endpoint_ids` - (Optional, List of Integer) IDs of configured notification endpoints.

### Alert Configuration

The `alert_configuration` block supports both log alert and metric alert fields.

* `type` - (Required, String, ForceNew) Alert type. Must be `LOG_ALERT` or `METRIC_ALERT`. Changing this forces a new resource.
* `suppress_notifications_minutes` - (Optional, Integer) Mute period after alert fires (log alerts).
* `alert_output_template_type` - (Optional, String) Output format for log alerts. Must be `JSON` or `TEXT`.
* `search_timeframe_minutes` - (Optional, Integer) Time window in minutes for log evaluation.
* `severity` - (Optional, String) Alert severity for metric alerts. Valid values: `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
* `sub_components` - (Optional, List of Block) Detection rules for log alerts. See [Sub Component](#sub-component) below.
* `correlations` - (Optional, Block) Correlation logic between sub-components (log alerts). See [Correlations](#correlations) below.
* `schedule` - (Optional, Block) Cron-based evaluation schedule (log alerts). See [Schedule](#schedule) below.
* `trigger` - (Optional, Block) Trigger configuration for metric alerts. See [Metric Trigger](#metric-trigger) below.
* `queries` - (Optional, List of Block) Metric queries for metric alerts. See [Metric Query](#metric-query) below.

#### Sub Component

The `sub_components` block supports:

* `query_definition` - (Required, Block) The query configuration. See [Query Definition](#query-definition) below.
* `trigger` - (Required, Block) Trigger conditions. See [Sub Component Trigger](#sub-component-trigger) below.
* `output` - (Optional, Block) Output configuration. See [Sub Component Output](#sub-component-output) below.

#### Query Definition

The `query_definition` block supports:

* `query` - (Required, String) Lucene/Elasticsearch query string (e.g., `"level:ERROR AND service:checkout"`).
* `filters` - (Optional, String) Boolean filters as JSON string. Example shape:

```json
{
  "bool": {
    "must": [],
    "should": [],
    "filter": [],
    "must_not": []
  }
}
```
* `group_by` - (Optional, List of String) Fields to group results by.
* `aggregation` - (Optional, Block) How to aggregate matching logs. See [Aggregation](#aggregation) below.
* `should_query_on_all_accounts` - (Optional, Boolean) Whether to query all accessible accounts. Default: `true`.
* `account_ids_to_query_on` - (Optional, List of Integer) Required if `should_query_on_all_accounts = false`.

#### Aggregation

The `aggregation` block supports:

* `aggregation_type` - (Required, String) Type of aggregation. Valid values: `SUM`, `MIN`, `MAX`, `AVG`, `COUNT`, `UNIQUE_COUNT`, `NONE`, `PERCENTAGE`, `PERCENTILE`.
* `field_to_aggregate_on` - (Optional, String) Field to aggregate on.
* `value_to_aggregate_on` - (Optional, String) Value to aggregate on.

#### Sub Component Trigger

The `trigger` block supports:

* `operator` - (Required, String) Comparison operator. Valid values: `LESS_THAN`, `GREATER_THAN`, `LESS_THAN_OR_EQUALS`, `GREATER_THAN_OR_EQUALS`, `EQUALS`, `NOT_EQUALS`.
* `severity_threshold_tiers` - (Required, List of Block) Severity thresholds. At least one required. See [Severity Threshold Tier](#severity-threshold-tier) below.

#### Severity Threshold Tier

The `severity_threshold_tiers` block supports:

* `severity` - (Required, String) Severity level. Valid values: `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
* `threshold` - (Required, Float) Threshold value.

**Important:** Threshold ordering depends on the operator:
- **For `GREATER_THAN`/`GREATER_THAN_OR_EQUALS`:** Higher severity must have higher thresholds (e.g., HIGH: 100, MEDIUM: 50, LOW: 10)
- **For `LESS_THAN`/`LESS_THAN_OR_EQUALS`:** Higher severity must have lower thresholds (e.g., HIGH: 10, MEDIUM: 50, LOW: 100)
- **For `EQUALS`/`NOT_EQUALS`:** Standard ordering applies

#### Sub Component Output

The `output` block supports:

* `should_use_all_fields` - (Optional, Boolean) Whether to use all fields in output. Default: `false`.
* `columns` - (Optional, List of Block) Column configurations. See [Column Config](#column-config) below.

**Important:** Custom `columns` are **only valid when `aggregation_type = "NONE"`**.

- If using any aggregation (`COUNT`, `SUM`, `AVG`, `MIN`, `MAX`, `UNIQUE_COUNT`): **Must set** `should_use_all_fields = true` and **cannot specify** `columns`.
- If using `aggregation_type = "NONE"`: Can set `should_use_all_fields = false` and specify custom `columns`.

#### Column Config

The `columns` block supports:

* `field_name` - (Required, String) Field name.
* `regex` - (Optional, String) Regular expression for field extraction.
* `sort` - (Optional, String) Sort direction. Valid values: `ASC`, `DESC`.

#### Schedule

The `schedule` block supports:

* `cron_expression` - (Required, String) Standard cron expression (e.g., `"*/5 * * * *"` = every 5 minutes).
* `timezone` - (Optional, String) Timezone for the cron expression. Default: `UTC`.

#### Correlations

The `correlations` block supports:

* `correlation_operators` - (Optional, List of String) Correlation operators (e.g., `["AND"]`).
* `joins` - (Optional, List of Map) Join configurations.

### Metric Trigger

The `trigger` block (inside `alert_configuration`, for metric alerts) supports:

* `type` - (Required, String) Trigger type. Valid values: `threshold`, `math`.
* `condition` - (Optional, Block) Threshold condition. Required when `type = "threshold"`. See [Trigger Condition](#trigger-condition) below.
* `expression` - (Optional, String) Math expression. Required when `type = "math"`. Uses query ref_ids (e.g., `"(A / B) * 100"`).

#### Trigger Condition

The `condition` block supports:

* `operator_type` - (Required, String) Comparison operator. Valid values: `above`, `below`, `within_range`, `outside_range`.
* `threshold` - (Optional, Float) Threshold value. Used with `above` and `below`.
* `from` - (Optional, Float) Range start. Used with `within_range` and `outside_range`.
* `to` - (Optional, Float) Range end. Used with `within_range` and `outside_range`.

### Metric Query

The `queries` block supports:

* `ref_id` - (Required, String) Query identifier (e.g., "A", "B") for use in math expressions.
* `query_definition` - (Required, Block) The query configuration. See [Metric Query Definition](#metric-query-definition) below.

#### Metric Query Definition

The `query_definition` block supports:

* `account_id` - (Optional, Integer) The account ID for the metrics data source.
* `promql_query` - (Required, String) PromQL query string (e.g., `"rate(http_requests_total[5m])"`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alert_id` - The unique alert identifier assigned by Logz.io.
* `created_at` - Unix timestamp (float) of alert creation.
* `updated_at` - Unix timestamp (float) of last update.
* `created_by` - Email of the user who created the alert.
* `updated_by` - Email of the user who last updated the alert.

## Import

Unified alerts can be imported using the alert type and ID, separated by a colon:

```bash
$ terraform import logzio_unified_alert.my_log_alert logs:alert-id-here
$ terraform import logzio_unified_alert.my_metric_alert metrics:alert-id-here
```

**Note:** When importing, you must specify both the alert type (`logs` or `metrics`) and the alert ID.
