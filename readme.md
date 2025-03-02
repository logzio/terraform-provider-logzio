# Logz.io Terraform provider

The Terraform Logz.io provider offers a great way to build integrations using Logz.io APIs.

Terraform is an infrastructure orchestrator written in Hashicorp Language (HCL). It is a popular Infrastructure-as-Code (IaC) tool that does away with manual configuration processes. You can take advantage of the Terraform Logz.io provider to really streamline the process of integrating observability into your dev workflows.

This guide assumes working knowledge of HashiCorp Terraform. If you're new to Terraform, we've got a great [introduction](https://logz.io/blog/terraform-vs-ansible-vs-puppet/) if you're in for one. We also recommend the official [Terraform guides and tutorials](https://www.terraform.io/guides/index.html).

### Capabilities

You can use the Terraform Logz.io Provider to manage users and log accounts in Logz.io, create and update log-based alerts and notification channels, and more.

The following Logz.io API endpoints are supported by this provider:

* [User management](https://api-docs.logz.io/docs/logz/manage-users)
* [Notification channels](https://api-docs.logz.io/docs/logz/manage-notification-endpoints)
* [Log-based alerts](https://github.com/logzio/public-api/tree/master/alerts)
* [Sub accounts](https://api-docs.logz.io/docs/logz/manage-time-based-log-accounts)
* [Alerts(v2)](https://api-docs.logz.io/docs/logz/alerts)
* [Log shipping token](https://api-docs.logz.io/docs/logz/manage-log-shipping-tokens)
* [Drop filters](https://api-docs.logz.io/docs/logz/drop-filters)
* [Archive logs](https://api-docs.logz.io/docs/logz/archive-logs)
* [Restore logs](https://api-docs.logz.io/docs/logz/restore-logs)
* [Authentication groups](https://api-docs.logz.io/docs/logz/authentication-groups)
* [Kibana objects](https://api-docs.logz.io/docs/logz/import-or-export-kibana-objects)
* [S3 Fetcher](https://api-docs.logz.io/docs/logz/connect-to-s-3-buckets)
* [Grafana Dashboards](https://api-docs.logz.io/docs/logz/create-dashboard)
* [Grafana folders](https://api-docs.logz.io/docs/logz/get-all-folders)
* [Grafana Alert Rules](https://api-docs.logz.io/docs/logz/get-alert-rules)
* [Grafana Contact Point](https://api-docs.logz.io/docs/logz/route-get-contactpoints)
* [Grafana Notification Policy](https://api-docs.logz.io/docs/logz/route-get-policy-tree)
* [Metrics Accounts](https://api-docs.logz.io/docs/logz/create-a-new-metrics-account)

#### Working with Terraform

<div class="tasklist">

**Before you begin, you'll need**:

* [Terraform CLI](https://learn.hashicorp.com/tutorials/terraform/install-cli)
* [Logz.io API token](/)

#### Get the Terraform Logz.io Provider

To install this provider, copy and paste this code into your Terraform configuration:

```hcl
terraform {
  required_providers {
    logzio = {
      source = "logzio/logzio"
    }
  }
}
```

This will install the latest Logz.io provider.
If you wish to use a specific version of the provider, add under `source` the field `version` and specify your preferred version.


##### Configuring the provider

The provider accepts the following arguments:

* **api_token** - (Required) The API token is used for authentication. [Learn more](/user-guide/tokens/api-tokens.html).

* **region** - (Defaults to null) The 2-letter region code identifies where your Logz.io account is hosted.
Defaults to null for accounts hosted in the US East - Northern Virginia region. [Learn more](https://docs.logz.io/user-guide/accounts/account-region.html)

###### Example

You can pass the variables in a bash command for the arguments:

```bash
provider "logzio" {
  api_token = var.api_token
  region= var.your_api_region
}
```
</div>


### Example - Create a new alert and a new Slack notification endpoint

Here's a great example demonstrating how easy it is to get up and running quickly with the Terraform Logz.io Provider.

This example adds a new Slack notification channel and creates a new alert in Kibana that will send notifications to the newly-created Slack channel.

The alert in this example will trigger whenever Logz.io records 10 loglevel:ERROR messages in 10 minutes.

```
terraform {
  required_providers {
    logzio = {
      source = "logzio/logzio"
    }
  }
}

provider "logzio" {
  api_token = "my-token"
  region= "au"
}

resource "logzio_endpoint" "my_endpoint" {
  title = "my_endpoint"
  description = "my slack endpoint"
  endpoint_type = "slack"
  slack {
    url = "${var.slack_url}"
  }
}

resource "logzio_alert_v2" "my_alert" {
  depends_on = [logzio_endpoint.my_endpoint]
  title = "hello_there"
  search_timeframe_minutes = 5
  is_enabled = false
  tags = ["some", "words"]
  suppress_notifications_minutes = 5
  alert_notification_endpoints = ["${logzio_endpoint.my_endpoint.id}"]
  output_type = "JSON"
  sub_components {
    query_string = "loglevel:ERROR"
    should_query_on_all_accounts = true
    operation = "GREATER_THAN"
    value_aggregation_type = "COUNT"
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

### Example - Create user

This example will create a user in your Logz.io account.

```
terraform {
  required_providers {
    logzio = {
      source = "logzio/logzio"
    }
  }
}

variable "api_token" {
  type = "string"
  description = "Your logzio API token"
}

variable "account_id" {
  description = "The account ID where the new user will be created"
}

provider "logzio" {
  api_token = var.api_token
  region = var.region
}

resource "logzio_user" "my_user" {
  username = "user_name@fun.io"
  fullname = "John Doe"
  roles = [ 2 ]
  account_id = var.account_id
}
```

Run the above plan using the following bash script:

```
export TF_LOG=DEBUG
terraform init
TF_VAR_api_token=${LOGZIO_API_TOKEN} TF_VAR_region=${LOGZIO_REGION} terraform plan -out terraform.plan
terraform apply terraform.plan
```

Before you run the script, update the arguments to match your details.
See our [examples](https://github.com/logzio/logzio_terraform_provider/tree/master/examples) for some complete working examples. 

### Contribute

Found a bug or want to suggest a feature? [Open an issue](https://github.com/logzio/logzio_terraform_provider/issues/new) about it.
Want to do it yourself? We are more than happy to accept external contributions from the community.
Simply fork the repo, add your changes and [open a PR](https://github.com/logzio/logzio_terraform_provider/pulls).

### Import sub-accounts as resources 

You can import multiple sub-accounts as follows:

```
terraform import logzio_subaccount.my_subaccount <SUBACCOUNT-ID>
```

### Trademark Disclaimer

Terraform is a trademark of HashiCorp, Inc.

