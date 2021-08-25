# Endpoint Resource

Use this data source to access information about existing Logz.io notification endpoints and custom webhooks.

* Endpoints can be used to send alerts, automate reports, share logs and dashboards, directly from Logz.io. Learn more about endpoint integrations in the [Logz.io Docs](https://docs.logz.io/user-guide/integrations/endpoints.html)
* Learn more about available [APIs for managing Logz.io endpoints](https://docs.logz.io/api/#tag/Manage-notification-endpoints)

## Example Usage

```hcl
resource "logzio_endpoint" "my_endpoint" {
  title = "my_endpoint"
  description = "hello"
  endpoint_type = "Slack"
  slack {
    url = "https://this.is.com/some/url"
  }
}
```

## Argument Reference

* `id` - ID of the notification endpoint.

## Attribute Reference

* `endpoint_type` - Specifies the endpoint resource type: `custom`, `slack`, `pagerduty`, `bigpanda`, `datadog`, `victorops`, `opsgenie`, `servicenow`, `microsoftteams`. Use the appropriate parameters for your selected endpoint type.
* `title` - Name of the endpoint.
* `description` - Detailed description of the endpoint.
## Endpoints used

* [Get all endpoints](https://docs.logz.io/api/#operation/getAllEndpoints)
* [Get endpoint by ID](https://docs.logz.io/api/#operation/getEndpointById)

Logz.io integrates with:
* [Slack](https://docs.logz.io/api#operation/createSlack)
* [PagerDuty](https://docs.logz.io/api/#operation/createPagerDuty)
* [BigPanda](https://docs.logz.io/api/#operation/createBigPanda)
* [Datadog](https://docs.logz.io/api/#operation/createDataDog)
* [VictorOps](https://docs.logz.io/api/#operation/createVictorops)
* [Custom integration](https://docs.logz.io/api/#operation/createCustom)
* [OpsGenie](https://docs.logz.io/api/#operation/createOpsGenie).
* [ServiceNow](https://docs.logz.io/api/#operation/createServiceNow).
* [Microsoft Teams](https://docs.logz.io/api/#operation/createMicrosoftTeams).

