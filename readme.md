# Logz.io Terraform provider

The Terraform Logz.io provider offers a great way to build integrations using Logz.io APIs.

Terraform is an infrastructure orchestrator written in Hashicorp Language (HCL). It is a popular Infrastructure-as-Code (IaC) tool that does away with manual configuration processes. You can take advantage of the Terraform Logz.io provider to really streamline the process of integrating observability into your dev workflows.

This guide assumes working knowledge of HashiCorp Terraform. If you're new to Terraform, we've got a great [introduction](https://logz.io/blog/terraform-vs-ansible-vs-puppet/) if you're in for one. We also recommend the official [Terraform guides and tutorials](https://www.terraform.io/guides/index.html).

### Capabilities

You can use the Terraform Logz.io Provider to manage users and log accounts in Logz.io, create and update log-based alerts and notification channels, and more.

The following Logz.io API endpoints are supported by this provider:

* [User management](https://docs.logz.io/api/#tag/Manage-users)
* [Notification channels](https://docs.logz.io/api/#tag/Manage-notification-endpoints)
* [Log-based alerts](https://github.com/logzio/public-api/tree/master/alerts)
* [Sub accounts](https://docs.logz.io/api/#tag/Manage-sub-accounts)
* [Alerts(v2)](https://docs.logz.io/api/#tag/Alerts)
* [Log shipping token](https://docs.logz.io/api/#tag/Manage-log-shipping-tokens)
* [Drop filters](https://docs.logz.io/api/#tag/Drop-filters)
* [Archive logs](https://docs.logz.io/api/#tag/Archive-logs)
* [Restore logs](https://docs.logz.io/api/#tag/Restore-logs)
* [Authentication groups](https://docs.logz.io/api/#tag/Authentication-groups)
* [Kibana objects](https://docs.logz.io/api/#tag/Import-or-export-Kibana-objects)

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
  api_token = "8387abb8-4831-53af-91de-5cd3784d9774"
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

### Changelog

- **v1.10.0**:
  - **Breaking Changes**:
    - Upgrading `terraform-plugin-sdk` to `v2.21.0`:
        - **Terraform 0.11** and lower will not be supported.
        - To read more about migrating to v2 of the `terraform-plugin-sdk` see [this article](https://www.terraform.io/plugin/sdkv2/guides/v2-upgrade-guide).
    - Removal of the `logzio_alert` resource. Use `logzio_alert_v2` instead.
    - `logzio_archive_logs`:
      - Removal of `s3_secret_credentials`. Use `aws_access_key` and `aws_secret_key` instead. See documentation for current configuration structure.
      - Refactor code to match current API.

<details>
  <summary markdown="span"> Expand to check old versions </summary>

- **v1.9.2**
    - *Bug fix*: Fix diff for resource `alert_v2` in fields `alert_notification_endpoints`, `notification_emails` ([#116](https://github.com/logzio/terraform-provider-logzio/issues/116)).
- **v1.9.1**
    - *Bug fix*: plugin won't crash when import for `archive_logs` fails.
- **v1.9.0**
    - Update client version(v1.11.0).
    - Support [Kibana objects](https://docs.logz.io/api/#tag/Import-or-export-Kibana-objects)
- **v1.8.3**
    - Update client version(v1.10.3).
    - Bug fixes:
        - **alerts_v2**:
          - Fix noisy diff for tags.
          - Field `is_enabled` defaults to `true`.
        - **sub_accounts**: allow creating flexible account without `max_daily_gb`.
- **v1.8.2**
    - Update client version(v1.10.2).
    - Bug fixes:
      - **alerts_v2**: fix bug for columns requiring sort field.
      - **sub_accounts**: add backoff for creating and updating sub accounts.
- **v1.8.1**
    - Upgrade provider's Go version from 1.15 to 1.16 in code and in release workflow.
    - Improve tests - add sleep after each test.
- **v1.8.0**
    - **Breaking change**: **custom endpoint** - refactor Headers - now a string of comma-seperated key-value pairs.
    - Update client version (v1.10.1) - bug fix for empty Header field.
    - Add to custom endpoint datasource Description field.
- **v1.7.0**
  - Update client version (v1.10).
  - Support [authentication groups resource](https://docs.logz.io/api/#tag/Authentication-groups).
  - `alerts_v2`: fix noisy diff for `severity_threshold_tiers`.

- **v1.6.1**
    - Update client version (v1.9.1) - bug fix for not found messages.
- **v1.6**
  - Update client version (v1.9).
  - Support [archive logs resource](https://docs.logz.io/api/#tag/Archive-logs).
  - Support [restore logs resource](https://docs.logz.io/api/#tag/Restore-logs).
- **v1.5**
    - Update client version(v1.8).
    - `sub_account`:
        - **Breaking changes:**
            - Deprecated attribute `utilization_settings`. Use `frequency_minutes` and `utilization_enabled` instead. 
        - Added attributes `flexible` and `reserved_daily_gb`.
        - Refactor tests.
        - Refactor code to use Terraform's retry.
    - `endpoint`:
        - **Breaking changes:**
            - Types naming.
        - New endpoints for OpsGenie, ServiceNow and Microsoft Teams.
        - Fix bug for `body_template` for `custom` type ([#70](https://github.com/logzio/terraform-provider-logzio/issues/70)).
        - Refactor tests.
        - Refactor code to use Terraform's retry.
    - `alerts_v2`:
        - Fix bug for `filter_must`,`filter_must_not` ([#82](https://github.com/logzio/terraform-provider-logzio/issues/82))
        - Refactor tests.
        - Refactor code to use Terraform's retry.
    - `drop_filter`:
        - Improve tests.
- v1.4
    - Update client version(v1.7).
    - Support [Drop Filter](https://docs.logz.io/api/#tag/Drop-filters) resource.
- v1.3
    - Update client version(v1.6).
    - Support Log Shipping Token resource.
- v1.2.4
    - Update client version(v1.5.3).
    - Fix `sub account` to return attributes `account_token` & `account_id`.
- v1.2.3
    - Fix bug for `custom endpoint` empty headers.
    - Allow empty sharing accounts array in `sub account`.
    - Add retry in resource `sub account`.
    - Replace module `terraform` with `terraform-plugin-sdk`. See further explanation [here](https://www.terraform.io/docs/extend/guides/v1-upgrade-guide.html).
    - Upgrade to Go v1.15.
    - Update client version(v1.5.2).
- v1.2.2
    - Update client version(v1.5.1).
    - Fix alerts_v2 sort bug.
- v1.2.1
    - Fix alerts_v2 type "TABLE" bug.
- v1.2
    - Update client version(v1.5.0).
    - Support Alerts v2 resource.
    - Fix 404 error for Alerts.
- v1.1.8
    - Update client version 
    - Fix custom endpoint headers bug
- v1.1.7
    - Published to Terraform registry    
- v1.1.5
    - Fix boolean parameters not parsed bug
    - Support import command to state
- v1.1.4
    - Support Sub Accounts resource
    - few bug fixes
    - removed circleCI  
- v1.1.3 
    - examples now use TF12
    - will now generate the meta data needed for the IntelliJ type IDE HCL plugin
    - no more travis - just circle CI
    - version bump to use the latest TF library (0.12.6), now compatible with TF12
- 1.1.2 
    - Moved some of the source code around to comply with TF provider layout convention
    - Moved the examples into an examples directory
  
</details>

