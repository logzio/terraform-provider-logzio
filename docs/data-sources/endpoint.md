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

* `endpoint_type` - Specifies the endpoint resource type: `custom`, `slack`, `pager_duty`, `big_panda`, `data_dog`, `victorops`. Use the appropriate parameters for your selected endpoint type.
* `title` - Name of the endpoint.
* `description` - Detailed description of the endpoint.
## Endpoints used

Logz.io integrates with:
* [Slack](https://docs.logz.io/api/#operation/updateSlack)
* [PagerDuty](https://docs.logz.io/api/#operation/updatePagerDuty)
* [BigPanda](https://docs.logz.io/api/#operation/updateBigPanda)
* [Datadog](https://docs.logz.io/api/#operation/updateDataDog)
* [VictorOps](https://docs.logz.io/api/#operation/updateVictorops)
* [Custom integration](https://docs.logz.io/api/#operation/updateCustom)

Other endpoints:
* [Get all endpoints](https://docs.logz.io/api/#operation/getAllEndpoints)
* [Get endpoint by ID](https://docs.logz.io/api/#operation/getEndpointById)