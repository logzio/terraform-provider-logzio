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
}
