# Logz.io Terraform provider

### Supports CRUD of Logz.io users, alerts and notification endpoints

This provider is based on the Logz.io client library - https://github.com/jonboydell/logzio_client

#### Requirements
Terraform 0.10.x

#### Obtaining the provider

The easiest way to get the provider is to run the `./scripts/update_plugin.sh` and edit the `PROVIDER_VERSION` variable to get the desired provider version. To get the latest:
```bash
bash <(curl -s https://raw.githubusercontent.com/logzio/logzio_terraform_provider/master/scripts/update_plugin.sh) 
```

#### Using the provider

**Note**: We provide multiple API endpoints by region, find your API host in [accounts & regions](https://docs.logz.io/user-guide/accounts/account-region.html#regions-and-urls). Default: api.logz.io

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

#### Build from source

To build from the project root, this will copy it into your [plugins directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).  You can copy it into your Terraform templates folder too.
```bash
./scripts/build.sh
```
**Note**: This build would work on Unix system, for other OSs, set the `GOOS` env variable before running the build script. For example:
```bash
export GOOS=windows
```

#### Changelog?

- 1.1.3 - examples now use TF12
- 1.1.3 - will now generate the meta data needed for the IntelliJ type IDE HCL plugin
- 1.1.3 - no more travis - just circle CI
- 1.1.3 - version bump to use the latest TF library (0.12.6), now compatible with TF12
- 1.1.2 - Moved some of the source code around to comply with TF provider layout convention
- 1.1.2 - Moved the examples into an examples directory

#### Doens't work?

Do an [https://github.com/logzio/logzio_terraform_provider/issues](issue).
Or fix it yourself and do a [https://github.com/logzio/logzio_terraform_provider/pulls](PR).

#### License

[![FOSSA Status](https://app.fossa.io/api/projects/custom%2B8359%2Fgit%40github.com%3Ajonboydell%2Flogzio_terraform_provider.git.svg?type=large)](https://app.fossa.io/projects/custom%2B8359%2Fgit%40github.com%3Ajonboydell%2Flogzio_terraform_provider.git?ref=badge_large)
