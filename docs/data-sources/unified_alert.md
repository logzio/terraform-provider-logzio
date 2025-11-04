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

# Use the data source outputs
output "alert_title" {
  value = data.logzio_unified_alert.log_alert_by_id.title
}

output "alert_enabled" {
  value = data.logzio_unified_alert.log_alert_by_id.enabled
}
```

## Argument Reference

* `type` - (Required, String) Alert type. Must be either `LOG_ALERT` or `METRIC_ALERT`.
* `alert_id` - (Required, String) The unique alert identifier.

**Note:** Lookup by `title` is not currently supported. Use `alert_id` to retrieve specific alerts.

## Attributes Reference

See the [Unified Alert Resource](../resources/unified_alert.md) for details on all available attributes. All resource attributes are exported by this data source.

### Common Attributes

* `alert_id` - The unique alert identifier.
* `title` - Alert name.
* `type` - Alert type (`LOG_ALERT` or `METRIC_ALERT`).
* `description` - Alert description.
* `tags` - List of tags.
* `enabled` - Whether the alert is enabled.
* `created_at` - Unix timestamp of creation.
* `updated_at` - Unix timestamp of last update.
* `folder_id` - Folder UID.
* `dashboard_id` - Dashboard UID.
* `panel_id` - Panel ID.
* `runbook` - Runbook text.
* `rca` - Whether RCA is enabled.
* `rca_notification_endpoint_ids` - RCA notification endpoint IDs.
* `use_alert_notification_endpoints_for_rca` - Whether to use alert endpoints for RCA.

### Log Alert Attributes

When `type = "LOG_ALERT"`, the following nested block is available:

* `log_alert` - Log alert configuration with all nested attributes as documented in the resource.

### Metric Alert Attributes

When `type = "METRIC_ALERT"`, the following nested block is available:

* `metric_alert` - Metric alert configuration with all nested attributes as documented in the resource.

