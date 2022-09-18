# Endpoint Provider

Provides a Logz.io notification endpoint or custom webhook resource. This can be used to create and manage Logz.io endpoint integrations.

* Endpoints can be used to send alerts, automate reports, share logs and dashboards, directly from Logz.io. Learn more about endpoint integrations in the [Logz.io Docs](https://docs.logz.io/user-guide/integrations/endpoints.html)
* Learn more about available [APIs for managing Logz.io endpoints](https://docs.logz.io/api/#tag/Manage-notification-endpoints)

## Example Usage

```hcl
resource "logzio_endpoint" "my_endpoint" {
  title = "my_endpoint"
  description = "hello"
  endpoint_type = "slack"
  slack {
    url = "https://this.is.com/some/url"
  }
}
```


## Argument Reference

* `endpoint_type` - (Required) Specifies the endpoint resource type: `custom`, `slack`, `pagerduty`, `bigpanda`, `datadog`, `victorops`, `opsgenie`, `servicenow`, `microsoftteams`. Use the appropriate parameters for your selected endpoint type.
* `title` - (Required) Name of the endpoint.
* `description` - (Required) Detailed description of the endpoint.
* `slack` - (Optional) Relevant when `endpoint_type` is `slack`. Manages a webhook to a specific Slack channel.
	  * `url` - Slack webhook URL to a specific Slack channel.
* `pagerduty` - (Optional) Relevant when `endpoint_type` is `pagerduty`. Manages a webhook to PagerDuty.
	* `service_key` - API key generated from PagerDuty for the purpose of the integration.
* `bigpanda` - (Optional) Relevant when `endpoint_type` is `bigpanda`. Manages a webhook to BigPanda.
	  * `api_token` - API authentication token from BigPanda.
  	* `app_key` - Application key from BigPanda.
* `datadog` - (Optional) Relevant when `endpoint_type` is `datadog`. Manages a webhook to Datadog.
  	* `api_key` - API key from Datadog.
* `victorops` - (Optional) Relevant when `endpoint_type` is `victorops`. Manages a webhook to VictorOps.
  * `routing_key` - Alert routing key from VictorOps.
  * `message_type` - VictorOps REST API `message_type`.
  * `service_api_key` - API key from VictorOps.
* `custom` - (Optional) Relevant when `endpoint_type` is `custom`. Manages a custom webhook for your integration of choice.
    * `url` - Specifies the URL destination.
    * `method` - Selects the HTTP request method.
    * `headers` - Header parameters for the request. String, sent as comma-separated key-value pairs.
    * `body_template` - string of JSON object that serves as the template for the message body.
* `opsgenie` - (Optional) Relevant when `endpoint_type` is `opsgenie`. Manages a webhook to OpsGenie.
    * `api_key` - API key from OpsGenie, see https://docs.opsgenie.com/docs/logz-io-integration.
* `servicenow` - (Optional) Relevant when `endpoint_type` is `servicenow`. Manages a webhook to ServiceNow.
    * `username` - ServiceNow user name.
    * `password` - ServiceNow password.
    * `url` - Provide your instance URL to connect to your existing ServiceNow instance, i.e. https://xxxxxxxxx.service-now.com/.
* `microsoftteams` - (Optional) Relevant when `endpoint_type` is `microsoftteams`. Manages a webhook to Microsoft Teams.
    * `url` - Your Microsoft Teams webhook URL, see https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook.

## Attribute Reference

* `id` - ID of the notification endpoint.

### Import endpoint as resource

You can import endpoint as follows:

```
terraform import logzio_endpoint.my_endpoint <ENDPOINT-ID>
```

## Endpoints used

Logz.io integrates with:
* [Slack](https://docs.logz.io/api/#operation/createSlack)
* [PagerDuty](https://docs.logz.io/api/#operation/createPagerDuty)
* [BigPanda](https://docs.logz.io/api/#operation/createBigPanda)
* [Datadog](https://docs.logz.io/api/#operation/createDataDog)
* [VictorOps](https://docs.logz.io/api/#operation/createVictorops)
* [Custom integration](https://docs.logz.io/api/#operation/createCustom)
* [OpsGenie](https://docs.logz.io/api/#operation/createOpsGenie).
* [ServiceNow](https://docs.logz.io/api/#operation/createServiceNow).
* [Microsoft Teams](https://docs.logz.io/api/#operation/createMicrosoftTeams).

Other endpoints:
* [Get all endpoints](https://docs.logz.io/api/#operation/getAllEndpoints)
* [Get endpoint by ID](https://docs.logz.io/api/#operation/getEndpointById)