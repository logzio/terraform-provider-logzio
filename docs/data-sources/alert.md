# Alert Resource

Use this data source to access information about existing Logz.io Log Monitoring alerts.

* Learn more about log alerts in the [Logz.io Docs](https://docs.logz.io/user-guide/alerts/)

## Argument Reference

*	`title` - Alert title.
*	`id` - Logz.io alert ID.
## Attribute Reference

*	`created_at` - Date and time in UTC when the alert was first created.
*	`created_by` - Email of the user who first created the alert.
*	`search_timeframe_minutes` -  The time frame for evaluating the log data is a sliding window, with 1 minute granularity.
*	`operation` - Specifies the operator for evaluating the results. Enum: `LESS_THAN`, `GREATER_THAN`, `LESS_THAN_OR_EQUALS`, `GREATER_THAN_OR_EQUALS`, `EQUALS`, `NOT_EQUALS`.
*	`severity_threshold_tiers` - Set as many as 5 thresholds, each with its own severity level.
  *	`severity` - Defaults to `MEDIUM`. Labels the event with a severity tag. Available severity tags are: `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
  *	`threshold` - Number of logs per search timeframe.
*	`alert_notification_endpoints` - Add email addresses and/or endpoint channels to automatically receive notifications with sample data when the alert triggers.
* `notification_emails` - Add email addresses to automatically receive notifications with sample data when the alert triggers.
*	`description` - A description of the event, its significance, and suggested next steps or instructions for the team.
*	`query_string` - Search query in Lucene syntax. You can combine filters and a search query to specify the logs you are looking for. You can combine filters and a search query to specify the logs you are looking for.
*	`filter` - You can use `must` and `must_not` filters. Filters are more efficient compared to a query, so it's recommended to opt for a filter over a `query_string`, where possible.
*	`tags` - Tags for filtering alerts and triggered alerts. Can be used in Kibana Discover, Kibana dashboards, and more.
*	`group_by_aggregation_fields` - Specify 1-3 fields by which to group the results and count them. If you apply a group by operation, the alert returns a count of the results aggregated by unique values.
*	`is_enabled` - True by default. If `true`, the alert is currently active.
*	`last_triggered_at` - Date and time in UTC when the alert last triggered.
*	`last_updated` - Date and time in UTC when the alert was last updated.
*	`notification_emails` - Array of email addresses to be notified when the alert triggers.
*	`suppress_notifications_minutes` - Add a waiting period in minutes to space out notifications. (The alert will still trigger but will not send out notifications during the waiting period.)
*	`value_aggregation_field` - Specify the field on which to run the aggregation for the trigger condition.
* `value_aggregation_type` - Specifies the aggregation operator. Can be: `SUM`, `MIN`, `MAX`, `AVG`, `COUNT`, `UNIQUE_COUNT`, `NONE`.