# Alert Provider (V2)

Provides a Logz.io log monitoring alert resource. This can be used to create and manage Logz.io log monitoring alerts.

* Learn more about log alerts in the [Logz.io Docs](https://docs.logz.io/user-guide/alerts/)

## Example Usage

```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_alert_v2" "my_alert" {
  title = "my_other_title"
  search_timeframe_minutes = 5
  is_enabled = false
  tags = ["some", "words"]
  alert_notification_endpoints = [logzio_endpoint.my_endpoint.id]
  suppress_notifications_minutes = 5
  output_type = "JSON"
  sub_components {
    query_string = "loglevel:ERROR"
    should_query_on_all_accounts = true
    operation = "GREATER_THAN"
    value_aggregation_type = "NONE"
    severity_threshold_tiers {
      severity = "HIGH"
      threshold = 10
    }
    severity_threshold_tiers {
      severity = "INFO"
      threshold = 5
    }
  }
}

```

## Argument Reference

### Required:
* `title` - (String) Alert title.
* `search_timeframe_minutes` - (Integer) The time frame for evaluating the log data is a sliding window, with 1 minute granularity.
* `sub_components` - (Block list) Sets the search criteria using a search query, filters, group by aggregations, accounts to search, and trigger conditions. See below for **nested schema**.

### Optional: 
* `description` - (String) A description of the event, its significance, and suggested next steps or instructions for the team.
* `tags` - (String list) Tags for filtering alerts and triggered alerts. Can be used in Kibana Discover, dashboards, and more.
* `is_enabled` - (Boolean) True by default. If `true`, the alert is currently active.
* `notification_emails` - (String list) Array of email addresses to be notified when the alert triggers.
* `alert_notification_endpoints` - (Integer list) Array of IDs of pre-configured endpoint channels to notify when the alert triggers.
* `suppress_notifications_minutes` - (Integer) Defaults to 5. Add a waiting period in minutes to space out notifications. (The alert will still trigger but will not send out notifications during the waiting period.)
* `output_type` - (String) Selects the output format for the alert notification. Can be: `"JSON"` or `"TABLE""` If the alert has no aggregations/group by fields, JSON offers the option to send full sample logs without selecting specific fields.
* `correlation_operator` - (String) Comma separated string of supported operators. Only applicable when multiple sub-components are in use. Selects a logic for correlating the alert’s sub-components. `AND` is currently the only supported operator. When AND is the correlation_operator, both sub-components must meet their triggering criteria for the alert to trigger.
* `joins` - (Map list) Specifies which group by fields must have the same values to trigger the alert. Joins the group by fields from the first and second sub-components. The key represents the index of the sub component in the array. The fields must be ordered pairs of the group by fields already in use in the `sub_components.query_string`.

#### Nested schema for `sub_components`:

##### Required:
* `query_string` - (String) Provide a Kibana search query written in Lucene syntax. The search query together with the filters select for the relevant logs. Cannot be null - send an asterisk wildcard `"*"` if not using a search query.
* `value_aggregation_type` - (String) Specifies the aggregation operator. Can be: `"SUM"`, `"MIN"`, `"MAX"`, `"AVG"`, `"COUNT"`, `"UNIQUE_COUNT"`, `"NONE"`. If `"COUNT"` or `"NONE"`, `value_aggregation_field` must be null, and `group_by_aggregation_fields` fields must not be empty. If any other operator type (other than `"NONE"` or `"COUNT"`), `value_aggregation_field` must not be null.
* `severity_threshold_tiers` - (Block) Sets a severity label per trigger threshold. If using more than one sub-component, only 1 severityThresholdTiers is allowed. Otherwise, 1 per enum are allowed (for a total of 5 thresholds of increasing severities). Increasing severity must adhere to the logic of the operator. See  below for **nested schema**.

##### Optional:
* `filter_must`(String) Runs Elasticsearch Bool Query filters on the data (before the search query is applied). The most efficient way to grab the logs you are looking for.
* `filter_must_not` - (String) Runs Elasticsearch Bool Query filters on the data (before the search query is applied). The most efficient way to grab the logs you are looking for.
* `group_by_aggregation_fields` - (String list) Specify 1-3 fields by which to group the results and count them. If you apply a group by operation, the alert returns a count of the results aggregated by unique values.
* `value_aggregation_field` - (String) Selects the field on which to run the aggregation for the trigger condition. Cannot be a field already in use for `group_by_aggregation_fields`.
* `should_query_on_all_accounts` - (Boolean) Defaults to true. Only applicable when the alert is run from the main account. If true, the alert runs on the main account and all associated searchable sub accounts. If false, specify relevant account IDs for the alert to monitor using the `account_ids_to_query_on` field.
* `account_ids_to_query_on` - (Integer list) Specify Account IDs to select which accounts the alert should monitor. The alert will be checked only on these accounts.
* `operation` - (String) Specifies the operator for evaluating the results. Can be: `"LESS_THAN"`, `"GREATER_THAN"`, `"LESS_THAN_OR_EQUALS"`, `"GREATER_THAN_OR_EQUALS"`, `"EQUALS"`, `"NOT_EQUALS"`.
* `columns` - (Block list) See  below for **nested schema**.

#### Nested schema for `sub_components.severity_threshold_tiers`:

##### Required:
* `severity` - (String) Labels the event with a severity tag. Available severity tags are: `"INFO"`, `"LOW"`, `"MEDIUM"`, `"HIGH"`, `"SEVERE"`.
* `threshold` - (Integer) Number of logs per search timeframe.

#### Nested schema for `sub_components.columns`:

##### Optional:
* `field_name` - (String) Specify the fields to be included in the notification.
* `regex` - (String) Trims the data using regex filters. [Learn more](https://docs.logz.io/user-guide/alerts/regex-filters.html).
* `sort` - (String) Specify a single field to sort by. The field cannot be an analyzed field (a field that supports free text search or searching by part of a message, such as the 'message' field). Should be: `"DESC"`, `"ASC"`.