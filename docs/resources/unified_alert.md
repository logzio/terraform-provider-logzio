# Unified Alert Resource

Provides a Logz.io unified alert resource. This resource allows you to create and manage both log-based and metric-based alerts through a single unified API.

**Note:** This is a POC (Proof of Concept) endpoint for the unified alerts API.

## Example Usage - Log Alert

```hcl
resource "logzio_unified_alert" "log_alert_example" {
  title       = "High Error Rate in Production"
  type        = "LOG_ALERT"
  description = "Triggers when error logs exceed threshold"
  tags        = ["production", "errors", "critical"]
  enabled     = true

  # Optional: RCA configuration
  rca                                     = true
  rca_notification_endpoint_ids           = [123]
  use_alert_notification_endpoints_for_rca = false

  # Optional: Dashboard linking
  folder_id    = "folder-uid"
  dashboard_id = "dashboard-uid"
  panel_id     = "panel-123"

  # Optional: Runbook
  runbook = "Check application logs and recent deployments"

  log_alert {
    search_timeframe_minutes = 15

    output {
      type                          = "JSON"
      suppress_notifications_minutes = 30

      recipients {
        emails                    = ["oncall@company.com"]
        notification_endpoint_ids = [456]
      }
    }

    sub_components {
      query_definition {
        query                        = "level:ERROR AND environment:production"
        should_query_on_all_accounts = true

        aggregation {
          aggregation_type = "COUNT"
        }
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
        should_use_all_fields = true
      }
    }

    # Optional: Schedule
    schedule {
      cron_expression = "*/5 * * * *"
      timezone        = "America/New_York"
    }
  }
}
```

## Example Usage - Metric Alert

```hcl
resource "logzio_unified_alert" "metric_alert_example" {
  title       = "High CPU Usage"
  type        = "METRIC_ALERT"
  description = "Fires when CPU usage exceeds 80%"
  tags        = ["infrastructure", "cpu", "critical"]
  enabled     = true

  metric_alert {
    severity = "HIGH"

    trigger {
      trigger_type            = "THRESHOLD"
      metric_operator         = "ABOVE"
      min_threshold           = 80.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = "prometheus-datasource-uid"
        promql_query   = "avg(rate(cpu_usage_total[5m])) * 100"
      }
    }

    recipients {
      emails                    = ["devops@company.com"]
      notification_endpoint_ids = [789]
    }
  }
}
```

## Example Usage - Metric Alert with Math Expression

```hcl
resource "logzio_unified_alert" "metric_math_alert" {
  title       = "5xx Error Rate Percentage High"
  type        = "METRIC_ALERT"
  description = "Fires when 5xx rate exceeds 2% of total requests"
  tags        = ["http", "errors"]
  enabled     = true

  metric_alert {
    severity = "HIGH"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "($A / $B) * 100"
      metric_operator          = "ABOVE"
      min_threshold            = 2.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"
      query_definition {
        datasource_uid = "prometheus-uid"
        promql_query   = "sum(rate(http_requests_total{status=~\"5..\"}[5m]))"
      }
    }

    queries {
      ref_id = "B"
      query_definition {
        datasource_uid = "prometheus-uid"
        promql_query   = "sum(rate(http_requests_total[5m]))"
      }
    }

    recipients {
      emails = ["alerts@company.com"]
    }
  }
}
```

## Argument Reference

### Required Arguments

* `title` - (String) Alert name displayed in UI and notifications.
* `type` - (String) Alert type. Must be either `LOG_ALERT` or `METRIC_ALERT`.

### Optional Arguments

* `description` - (String) Explanation of alert purpose and firing conditions.
* `tags` - (List of String) Labels for grouping and searching alerts.
* `folder_id` - (String) UID of the unified folder in Logz.io.
* `dashboard_id` - (String) UID of the unified dashboard for context linking.
* `panel_id` - (String) Specific panel ID on the dashboard.
* `runbook` - (String) Operational instructions for responders.
* `enabled` - (Boolean) Alert activation status. Default: `true`.
* `rca` - (Boolean) Enable Root Cause Analysis. Default: `false`.
* `rca_notification_endpoint_ids` - (List of Integer) Notification endpoint IDs for RCA results.
* `use_alert_notification_endpoints_for_rca` - (Boolean) When true, RCA uses same endpoints as alert. Default: `false`.
* `log_alert` - (Block, Max: 1) Log alert configuration. Required when `type = "LOG_ALERT"`. See [Log Alert](#log-alert) below.
* `metric_alert` - (Block, Max: 1) Metric alert configuration. Required when `type = "METRIC_ALERT"`. See [Metric Alert](#metric-alert) below.

### Log Alert

The `log_alert` block supports:

* `search_timeframe_minutes` - (Required, Integer) Time window in minutes for log evaluation.
* `output` - (Required, Block) Notification configuration. See [Log Alert Output](#log-alert-output) below.
* `sub_components` - (Required, List of Block) Detection rules. At least one required. See [Sub Component](#sub-component) below.
* `correlations` - (Optional, Block) Correlation logic between sub-components. See [Correlations](#correlations) below.
* `schedule` - (Optional, Block) Cron-based evaluation schedule. See [Schedule](#schedule) below.

#### Log Alert Output

The `output` block supports:

* `type` - (Required, String) Output format. Must be `JSON` or `TABLE`.
* `suppress_notifications_minutes` - (Optional, Integer) Mute period after alert fires.
* `recipients` - (Required, Block) Who receives notifications. See [Recipients](#recipients) below.

#### Recipients

The `recipients` block supports:

* `emails` - (Optional, List of String) Email addresses for notifications.
* `notification_endpoint_ids` - (Optional, List of Integer) IDs of configured notification endpoints.

**Note:** At least one of `emails` or `notification_endpoint_ids` should be provided.

#### Sub Component

The `sub_components` block supports:

* `query_definition` - (Required, Block) The query configuration. See [Query Definition](#query-definition) below.
* `trigger` - (Required, Block) Trigger conditions. See [Sub Component Trigger](#sub-component-trigger) below.
* `output` - (Optional, Block) Output configuration. See [Sub Component Output](#sub-component-output) below.

#### Query Definition

The `query_definition` block supports:

* `query` - (Required, String) Lucene/Elasticsearch query string (e.g., `"level:ERROR AND service:checkout"`).
* `filters` - (Optional, String) Boolean filters as JSON string.
* `group_by` - (Optional, List of String) Fields to group results by.
* `aggregation` - (Optional, Block) How to aggregate matching logs. See [Aggregation](#aggregation) below.
* `should_query_on_all_accounts` - (Optional, Boolean) Whether to query all accessible accounts. Default: `true`.
* `account_ids_to_query_on` - (Optional, List of Integer) Required if `should_query_on_all_accounts = false`.

#### Aggregation

The `aggregation` block supports:

* `aggregation_type` - (Required, String) Type of aggregation. Valid values: `SUM`, `MIN`, `MAX`, `AVG`, `COUNT`, `UNIQUE_COUNT`, `NONE`.
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

Example for `LESS_THAN` (detecting low values):
```hcl
trigger {
  operator = "LESS_THAN"
  
  severity_threshold_tiers {
    severity  = "HIGH"
    threshold = 10    # More critical: value < 10
  }
  
  severity_threshold_tiers {
    severity  = "MEDIUM"
    threshold = 50    # Less critical: value < 50
  }
}
```

#### Sub Component Output

The `output` block supports:

* `should_use_all_fields` - (Optional, Boolean) Whether to use all fields in output. Default: `false`.
* `columns` - (Optional, List of Block) Column configurations. See [Column Config](#column-config) below.

**Important:** Custom `columns` are **only valid when `aggregation_type = "NONE"`**. 

- If using any aggregation (`COUNT`, `SUM`, `AVG`, `MIN`, `MAX`, `UNIQUE_COUNT`): **Must set** `should_use_all_fields = true` and **cannot specify** `columns`.
- If using `aggregation_type = "NONE"`: Can set `should_use_all_fields = false` and specify custom `columns`.

Example with aggregation (no custom columns):
```hcl
query_definition {
  aggregation {
    aggregation_type = "COUNT"
  }
}

output {
  should_use_all_fields = true  # Required with aggregation
  # Cannot specify columns here
}
```

Example with NONE aggregation (custom columns allowed):
```hcl
query_definition {
  aggregation {
    aggregation_type = "NONE"
  }
}

output {
  should_use_all_fields = false
  
  columns {
    field_name = "timestamp"
    sort       = "DESC"
  }
  
  columns {
    field_name = "message"
  }
}
```

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

### Metric Alert

The `metric_alert` block supports:

* `severity` - (Required, String) Alert severity level. Valid values: `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
* `trigger` - (Required, Block) Trigger configuration. See [Metric Trigger](#metric-trigger) below.
* `queries` - (Required, List of Block) Metric queries. At least one required. See [Metric Query](#metric-query) below.
* `recipients` - (Required, Block) Who receives notifications. See [Recipients](#recipients) above.

#### Metric Trigger

The `trigger` block supports:

* `trigger_type` - (Required, String) Trigger type. Valid values: `THRESHOLD`, `MATH`.
* `metric_operator` - (Optional, String) Required for threshold triggers. Valid values: `ABOVE`, `BELOW`, `WITHIN_RANGE`, `OUTSIDE_RANGE`.
* `min_threshold` - (Optional, Float) Minimum threshold value.
* `max_threshold` - (Optional, Float) Maximum threshold value (required for `WITHIN_RANGE` and `OUTSIDE_RANGE`).
* `math_expression` - (Optional, String) Required when `trigger_type = "MATH"`. Expression using query ref_ids (e.g., `"$A / $B * 100"`).
* `search_timeframe_minutes` - (Required, Integer) Evaluation time window in minutes.

#### Metric Query

The `queries` block supports:

* `ref_id` - (Required, String) Query identifier (e.g., "A", "B") for use in math expressions.
* `query_definition` - (Required, Block) The query configuration. See [Metric Query Definition](#metric-query-definition) below.

#### Metric Query Definition

The `query_definition` block supports:

* `datasource_uid` - (Required, String) UID of the Prometheus/metrics datasource in Logz.io.
* `promql_query` - (Required, String) PromQL query string (e.g., `"rate(http_requests_total[5m])"`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alert_id` - The unique alert identifier assigned by Logz.io.
* `created_at` - Unix timestamp (float) of alert creation.
* `updated_at` - Unix timestamp (float) of last update.

## Import

Unified alerts can be imported using the alert type and ID, separated by a colon:

```bash
$ terraform import logzio_unified_alert.my_log_alert LOG_ALERT:alert-id-here
$ terraform import logzio_unified_alert.my_metric_alert METRIC_ALERT:alert-id-here
```

**Note:** When importing, you must specify both the alert type (`LOG_ALERT` or `METRIC_ALERT`) and the alert ID.

