# Unified Alert Data Source

Use this data source to access information about an existing Logz.io unified alert.

## Example Usage

```hcl
# Lookup log alert by ID
data "logzio_unified_alert" "log_alert_by_id" {
  type     = "LOG_ALERT"
  alert_id = "alert-123"
}

# Lookup metric alert by ID
data "logzio_unified_alert" "metric_alert_by_id" {
  type     = "METRIC_ALERT"
  alert_id = "alert-456"
}

# Reference from a managed resource
data "logzio_unified_alert" "from_resource" {
  type     = "LOG_ALERT"
  alert_id = logzio_unified_alert.my_log_alert.alert_id

  depends_on = [logzio_unified_alert.my_log_alert]
}

# Use the data source outputs
output "alert_title" {
  value = data.logzio_unified_alert.log_alert_by_id.title
}

output "alert_enabled" {
  value = data.logzio_unified_alert.log_alert_by_id.enabled
}

output "metric_severity" {
  value = data.logzio_unified_alert.metric_alert_by_id.alert_configuration[0].severity
}
```

## Argument Reference

* `type` - (Required, String) Alert type. Must be `LOG_ALERT` or `METRIC_ALERT`.
* `alert_id` - (Required, String) The unique alert identifier.

## Attributes Reference

See the [Unified Alert Resource](../resources/unified_alert.md) for details on all available attributes. All resource attributes are exported by this data source.

### Common Attributes

* `alert_id` - The unique alert identifier.
* `title` - Alert name.
* `description` - Alert description.
* `tags` - Set of tags.
* `enabled` - Whether the alert is enabled.
* `created_at` - Unix timestamp of creation.
* `updated_at` - Unix timestamp of last update.
* `created_by` - Email of the user who created the alert.
* `updated_by` - Email of the user who last updated the alert.
* `linked_panel` - Dashboard panel link with `folder_id`, `dashboard_id`, `panel_id`.
* `runbook` - Runbook text.
* `rca` - Whether RCA is enabled.
* `rca_notification_endpoint_ids` - RCA notification endpoint IDs.
* `use_alert_notification_endpoints_for_rca` - Whether to use alert endpoints for RCA.
* `recipients` - Notification recipients with `emails` and `notification_endpoint_ids`.
* `alert_configuration` - The alert configuration block. See the resource documentation for full details.
