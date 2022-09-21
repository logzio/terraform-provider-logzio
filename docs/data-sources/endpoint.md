# Endpoint Datasource

Use this data source to access information about existing Logz.io notification endpoints and custom webhooks.

* Endpoints can be used to send alerts, automate reports, share logs and dashboards, directly from Logz.io. Learn more about endpoint integrations in the [Logz.io Docs](https://docs.logz.io/user-guide/integrations/endpoints.html)
* Learn more about available [APIs for managing Logz.io endpoints](https://docs.logz.io/api/#tag/Manage-notification-endpoints).

## Argument Reference

* `id` - ID of the notification endpoint.

## Attribute Reference

* `endpoint_type` - Specifies the endpoint resource type: `custom`, `slack`, `pagerduty`, `bigpanda`, `datadog`, `victorops`, `opsgenie`, `servicenow`, `microsoftteams`. Use the appropriate parameters for your selected endpoint type.
* `title` - Name of the endpoint.
* `description` - Detailed description of the endpoint.
