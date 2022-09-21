# Alert V2 Datasource

Provides a Logz.io log monitoring alert resource. This can be used to create and manage Logz.io log monitoring alerts.

* Learn more about log alerts in the [Logz.io Docs](https://docs.logz.io/user-guide/alerts/)

## Argument Reference

* `title` - Alert title.
* `id` - Logz.io alert ID.

## Attribute Reference

* `created_at` - Date and time in UTC when the alert was first created.
* `created_by` - Email of the user who first created the alert.
* `updated_at` - Date and time in UTC when the alert was last updated.
* `updated by` - Email of the user who last updated the alert.
* `description` - A description of the event, its significance, and suggested next steps or instructions for the team.
* `tags` - Tags for filtering alerts and triggered alerts. Can be used in Kibana Discover, dashboards, and more.
* `search_timeframe_minutes` - The time frame for evaluating the log data is a sliding window, with 1 minute granularity.
* `is_enabled` - If `true`, the alert is currently active.
* `notification_emails` - Array of email addresses to be notified when the alert triggers.
* `alert_notification_endpoints` - Array of IDs of pre-configured endpoint channels to notify when the alert triggers.
* `suppress_notifications_minutes` - Add a waiting period in minutes to space out notifications. (The alert will still trigger but will not send out notifications during the waiting period.)
* `output_type` - Selects the output format for the alert notification. Can be: `"JSON"` or `"TABLE""` If the alert has no aggregations/group by fields, JSON offers the option to send full sample logs without selecting specific fields.
* `correlation_operator` - Comma separated string of supported operators. Only applicable when multiple sub-components are in use. Selects a logic for correlating the alert’s sub-components. `AND` is currently the only supported operator. When AND is the correlation_operator, both sub-components must meet their triggering criteria for the alert to trigger.
* `joins` - Specifies which group by fields must have the same values to trigger the alert. Joins the group by fields from the first and second sub-components. The key represents the index of the sub component in the array. The fields must be ordered pairs of the group by fields already in use in the `sub_components.query_string`.
* `sub_components` - Sets the search criteria using a search query, filters, group by aggregations, accounts to search, and trigger conditions.
* `sub_components.query_string` - Provide a Kibana search query written in Lucene syntax. The search query together with the filters select for the relevant logs. Cannot be null - send an asterisk wildcard `"*"` if not using a search query.
* `sub_components.filter_must` - Runs Elasticsearch Bool Query `must` filters on the data (before the search query is applied). The most efficient way to grab the logs you are looking for.
* `sub_components.filter_must_not` - Runs Elasticsearch Bool Query `must_not` filters on the data (before the search query is applied). The most efficient way to grab the logs you are looking for.
* `sub_components.group_by_aggregation_fields` - Specify 1-3 fields by which to group the results and count them. If you apply a group by operation, the alert returns a count of the results aggregated by unique values.
* `sub_components.value_aggregation_type` - Specifies the aggregation operator. Can be: `"SUM"`, `"MIN"`, `"MAX"`, `"AVG"`, `"COUNT"`, `"UNIQUE_COUNT"`, `"NONE"`. If `"COUNT"` or `"NONE"`, `value_aggregation_field` must be null, and `group_by_aggregation_fields` fields must not be empty. If any other operator type (other than `"NONE"` or `"COUNT"`), `value_aggregation_field` must not be null.
* `sub_components.value_aggregation_field` - Selects the field on which to run the aggregation for the trigger condition. Cannot be a field already in use for `group_by_aggregation_fields`.
* `sub_components.should_query_on_all_accounts` - Defaults to true. Only applicable when the alert is run from the main account. If true, the alert runs on the main account and all associated searchable sub accounts. If false, specify relevant account IDs for the alert to monitor using the `account_ids_to_query_on` field.
* `sub_components.account_ids_to_query_on` - Specify Account IDs to select which accounts the alert should monitor. The alert will be checked only on these accounts.
* `sub_components.operation` - Specifies the operator for evaluating the results. Can be: `"LESS_THAN"`, `"GREATER_THAN"`, `"LESS_THAN_OR_EQUALS"`, `"GREATER_THAN_OR_EQUALS"`, `"EQUALS"`, `"NOT_EQUALS"`.
* `sub_components.severity_threshold_tiers` - Sets a severity label per trigger threshold. If using more than one sub-component, only 1 severityThresholdTiers is allowed. Otherwise, 1 per enum are allowed (for a total of 5 thresholds of increasing severities). Increasing severity must adhere to the logic of the operator.
* `sub_components.severity_threshold_tiers.severity` - Labels the event with a severity tag. Available severity tags are: `"INFO"`, `"LOW"`, `"MEDIUM"`, `"HIGH"`, `"SEVERE"`.
* `sub_components.severity_threshold_tiers.threshold` - Number of logs per search timeframe.
* `sub_components.columns.field_name` - Specify the fields to be included in the notification. 
* `sub_components.columns.regex` - Trims the data using regex filters. [Learn more](https://docs.logz.io/user-guide/alerts/regex-filters.html).
* `sub_components.columns.sort` - Specify a single field to sort by. The field cannot be an analyzed field (a field that supports free text search or searching by part of a message, such as the 'message' field). Should be: `"DESC"`, `"ASC"`.
* `sub_components.output_should_use_all_fields` - If true, the notification output will include entire logs with all of their fields in the sample data.
