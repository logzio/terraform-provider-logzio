# Logz.io Terraform provider

|branch|build status|
|---|---|
|develop|[![CircleCI](https://circleci.com/gh/jonboydell/logzio_terraform_provider/tree/develop.svg?style=svg)](https://circleci.com/gh/jonboydell/logzio_terraform_provider/tree/develop)|

### Supports CRUD of Logz.io user, alerts and notification endpoints

This provider is based on the Logz.io client library - https://github.com/jonboydell/logzio_client

#### What's new?

- 1.1.3 - examples now use TF12
- 1.1.3 - will now generate the meta data needed for the IntelliJ type IDE HCL plugin
- 1.1.3 - no more travis - just circle CI
- 1.1.3 - version bump to use the latest TF library (0.12.6), now compatible with TF12
- 1.1.2 - Moved some of the source code around to comply with TF provider layout convention
- 1.1.2 - Moved the examples into an examples directory

#### Obtaining the provider

The easiest way to get the provider and the JetBrains IDE HCL meta-data is to run the `./scripts/update_plugin.sh` and edit the `PROVIDER_VERSION` variable to get the right provider version.

However...

To build; from the project root (on a *nix style system), this will copy it into your [plugins directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).  You can copy it into your Terraform templates folder too.
```bash
./scripts/build.sh
```

#### Using the provider

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
  severity_threshold_tiers {
      severity = "HIGH",
      threshold = 10
  }
}
```

#### Running the tests
`GO111MODULE=on TF_ACC=true go test -v .`

#### Doens't work?

Do an [https://github.com/jonboydell/logzio_terraform_provider/issues](issue).
Fix it yourself and do a [https://github.com/jonboydell/logzio_terraform_provider/pulls](PR), please create any fix branches from `develop`.  They'll be merged back into `develop` and go `master` from there.  Releases are from `master`.

#### License

[![FOSSA Status](https://app.fossa.io/api/projects/custom%2B8359%2Fgit%40github.com%3Ajonboydell%2Flogzio_terraform_provider.git.svg?type=large)](https://app.fossa.io/projects/custom%2B8359%2Fgit%40github.com%3Ajonboydell%2Flogzio_terraform_provider.git?ref=badge_large)
