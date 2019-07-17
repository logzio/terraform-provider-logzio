# Logz.io Terraform provider

develop
[![Build Status](https://travis-ci.org/jonboydell/logzio_terraform_provider.svg?branch=develop)](https://travis-ci.org/jonboydell/logzio_terraform_provider)
[![Coverage Status](https://coveralls.io/repos/github/jonboydell/logzio_terraform_provider/badge.svg?branch=develop)](https://coveralls.io/github/jonboydell/logzio_terraform_provider?branch=develop)

master
[![Build Status](https://travis-ci.org/jonboydell/logzio_terraform_provider.svg?branch=master)](https://travis-ci.org/jonboydell/logzio_terraform_provider)
[![Coverage Status](https://coveralls.io/repos/github/jonboydell/logzio_terraform_provider/badge.svg?branch=master)](https://coveralls.io/github/jonboydell/logzio_terraform_provider?branch=master)

### Supports CRUD of Logz.io alerts and notification endpoints

This provider is based on the Logz.io client library - https://github.com/jonboydell/logzio_client

##### Obtaining the provider

To build; from the project root (on a *nix style system), this will copy it into your [plugins directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).  You can copy it into your Terraform templates folder too.
```bash
./scripts/build.sh
```

You can [get a release from here](https://github.com/jonboydell/logzio_terraform_provider/releases) and follow these [instructions](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
You'll need to do a `terraform init` for it to pick up the provider.

##### Using the provider

Note: logz.io provides multiple endpoints for their service, if you are not using the default, `https://api.logz.io` then you'll have to specify an override in the provider.
```hcl-terraform
provider "logzio" {
  api_token = "${var.api_token}"
  base_url = "${var.your_api_endpoint}" #e.g. https://api-au.logz.io
}
```

This simple example will create a Logz.io Slack notification endpoint (you'll need to provide the right URL) and an alert that
is triggered should Logz.io record 10 loglevel:ERROR messages in 5 minutes.  To make this example work you will also need to provide
your Logz.io API token.

```hcl-terraform
provider "logzio" {
  api_token = "${var.api_token}"
}

resource "logzio_endpoint" "my_endpoint" {
  title = "my_endpoint"
  description = "hello"
  endpoint_type = "slack"
  slack {
    url = "${var.slack_url}"
  }
}

resource "logzio_alert" "my_alert" {
  title = "my_other_title"
  query_string = "loglevel:ERROR"
  operation = "GREATER_THAN"
  notification_emails = []
  search_timeframe_minutes = 5
  value_aggregation_type = "NONE"
  alert_notification_endpoints = ["${logzio_endpoint.my_endpoint.id}"]
  suppress_notifications_minutes = 5
  severity_threshold_tiers = [
    {
      "severity" = "HIGH",
      "threshold" = 10
    }
  ]
}
```
## How to run the tests
1. `dep ensure -v`
2. `TF_ACC=true go test -v .`

##### Doens't work?

Do an [https://github.com/jonboydell/logzio_terraform_provider/issues](issue).
Fix it yourself and do a [https://github.com/jonboydell/logzio_terraform_provider/pulls](PR), please create any fix branches from `develop`.  They'll be merged back into `develop` and go `master` from there.  Releases are from `master`.

#### License

[![FOSSA Status](https://app.fossa.io/api/projects/custom%2B8359%2Fgit%40github.com%3Ajonboydell%2Flogzio_terraform_provider.git.svg?type=large)](https://app.fossa.io/projects/custom%2B8359%2Fgit%40github.com%3Ajonboydell%2Flogzio_terraform_provider.git?ref=badge_large)