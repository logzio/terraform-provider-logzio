# Grafana Contact Point

Provides a Logz.io Grafana Contact Point resource. This can be used to create and manage Logz.io Grafana Contact Points.

* Learn more about grafana contact points in the [Logz.io Docs](https://docs.logz.io/api/#tag/Grafana-contact-points).

## Example Usage

```hcl
variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_grafana_contact_point" "test_cp_email" {
  name = "my-email-cp"
  email {
    addresses = ["example@example.com", "example2@example.com"]
    disable_resolve_message = false
    single_email = true
    message = "{{ len .Alerts.Firing }} firing."
  }
}

resource "logzio_grafana_contact_point" "test_cp_googlechat" {
  name = "my-googlechat-cp"
  googlechat {
    url = "some.url"
    disable_resolve_message = false
    message = "{{ len .Alerts.Firing }} firing."
  }
}

resource "logzio_grafana_contact_point" "test_cp_opsgenie" {
  name = "my-opsgenie-cp"
  opsgenie {
    disable_resolve_message = false
    api_url = "some.url"
    api_key = "some_api_key"
    auto_close = false
    override_priority = true
    send_tags_as = "both"
  }
}

resource "logzio_grafana_contact_point" "test_cp_pagerduty" {
  name = "my-pagerduty-cp"
  pagerduty {
    integration_key = "some-key"
    class = "some-class"
    component = "some-component"
    group = "some-group"
    severity = "info"
    disable_resolve_message = false
  }
}

resource "logzio_grafana_contact_point" "test_cp_slack" {
  name = "my-slack-cp"
  slack {
    token = "some-token"
    title = "some-title"
    text = "{{ len .Alerts.Firing }} firing."
    mention_channel = "here"
    recipient = "me"
    disable_resolve_message = false
  }
}

resource "logzio_grafana_contact_point" "test_cp_teams" {
  name = "my-teams-cp"
  teams {
    url = "url"
    message = "message"
    disable_resolve_message = false
  }
}

resource "logzio_grafana_contact_point" "test_cp_victorops" {
  name = "my-victorops-cp"
  victorops {
    url = "url"
    message_type = "CRITICAL"
    disable_resolve_message = false
  }
}

resource "logzio_grafana_contact_point" "test_cp_webhook" {
  name = "my-webhook-cp"
  webhook {
    url = "url"
    disable_resolve_message = false
  }
  opsgenie {
    disable_resolve_message = false
    api_url = "some.url"
    api_key = "some_api_key"
    auto_close = false
    override_priority = true
    send_tags_as = "both"
  }
}

```

## Argument Reference

### Required:

* `name` - (String) The name of your contact point.

### Optional:

* `email` - (Block List) A contact point that sends notifications to an email address. See below for **nested schema**.
* `googlechat` - (Block List) A contact point that sends notifications to Google Chat. See below for **nested schema**.
* `opsgenie` - (Block List) A contact point that sends notifications to OpsGenie. See below for **nested schema**.
* `pagerduty` - (Block List) A contact point that sends notifications to PagerDuty. See below for **nested schema**.
* `slack` - (Block List) A contact point that sends notifications to Slack. See below for **nested schema**.
* `teams` - (Block List) A contact point that sends notifications to Microsoft Teams. See below for **nested schema**.
* `victorops` - (Block List) A contact point that sends notifications to VictorOps. See below for **nested schema**.
* `webhook` - (Block List) A contact point that sends notifications to an arbitrary webhook. See below for **nested schema**.

##  Attribute Reference

* `id` - (String) The ID of this resource.

## Nested schema for `email`:

### Argument Reference

#### Required:

* `addresses` - (List of String) The addresses to send emails to.

#### Optional:

* `single_email` - (Boolean) Whether to send a single email CC'ing all addresses, rather than a separate email to each address. Defaults to `false`.
* `message` - (String) The templated content of the email. Defaults to ``.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Nested schema for `googlechat`:

#### Required:

* `url` - (String, Sensitive) The Google Chat webhook URL.

#### Optional:

* `message` - (String) The templated content of the message.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Nested schema for `opsgenie`:

#### Required:

* `api_key` - (String, Sensitive) The OpsGenie API key to use.

#### Optional:

* `api_url` - (String) Allows customization of the OpsGenie API URL.
* `auto_close` - (Boolean) Whether to auto-close alerts in OpsGenie when they resolve in the Alertmanager.
* `override_priority` - (Boolean) Whether to allow the alert priority to be configured via the value of the og_priority annotation on the alert.
* `send_tags_as` - (String) Whether to send annotations to OpsGenie as Tags, Details, or both. Supported values are `tags`, `details`, `both`, or empty to use the default behavior of Tags.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Nested schema for `pagerduty`:

#### Required:

* `integration_key` - (String, Sensitive) The PagerDuty API key.

#### Optional:

* `class` - (String) The class or type of event.
* `component` - (String) The component being affected by the event.
* `group` - (String) The group to which the provided component belongs to.
* `summary` - (String) The templated summary message of the event.
* `severity` - (String) The PagerDuty event severity level. Can be one of `info`, `warning`, `error`, `critical`.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Nested schema for `slack`:

#### Required:

* `recipient` - (String) Channel, private group, or IM channel (can be an encoded ID or a name) to send messages to.

#### Optional:

* `endpoint_url` - (String) Use this to override the Slack API endpoint URL to send requests to.
* `mention_channel` - (String) Describes how to ping the slack channel that messages are being sent to. Options are `here` for an @here ping, `channel` for @channel, or empty for no ping.
* `mention_groups` - (String) Comma-separated list of groups to mention in the message.
* `mention_users` - (String) Comma-separated list of users to mention in the message.
* `text` - (String) Templated content of the message.
* `title` - (String) Templated title of the message.
* `token` - (String, Sensitive) A Slack API token,for sending messages directly without the webhook method.
* `url` - (String, Sensitive) A Slack webhook URL,for sending messages via the webhook method.
* `username` - (String) Username for the bot to use.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Nested schema for `teams`:

#### Required:

* `url` - (String, Sensitive) A Teams webhook URL.

#### Optional:

* `message` - (String) The templated message content to send.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Nested schema for `victorops`:

#### Required:

* `url` - (String) The VictorOps webhook URL.

#### Optional:

* `message_type` - (String) The VictorOps alert state - typically either `CRITICAL` or `WARNING`.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Nested schema for `webhook`:

#### Required:

* `url` - (String) The URL to send webhook requests to.

#### Optional:

* `http_method` - (String) The HTTP method to use in the request. Can be either `PUT` or `POST`.
* `max_alerts` - (Number) The maximum number of alerts to send in a single request. This can be helpful in limiting the size of the request body. The default is 0, which indicates no limit.
* `password` - (String, Sensitive) The username to use in basic auth headers attached to the request. If omitted, basic auth will not be used.
* `username` - (String) The username to use in basic auth headers attached to the request. If omitted, basic auth will not be used.
* `disable_resolve_message` - (Boolean) Whether to disable sending resolve messages. Defaults to `false`.
* `settings` - (Map of String, Sensitive) Additional custom properties to attach to the notifier. Defaults to `map[]`.

###  Attribute Reference

* `uid` - (String) The UID of the contact point.

## Import contact point as resource

You can import contact point as follows:

```
terraform import logzio_grafana_contact_point.my_cp <CONTACT-POINT-NAME>
```
