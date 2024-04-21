# Metrics accounts Provider

Provides a Logz.io Metrics account management resource. This can be used to create and manage Logz.io metric accounts.

* Learn more about accounts in the [Logz.io Docs](https://docs.logz.io/docs/user-guide/admin/logzio-accounts/manage-the-main-account-and-sub-accounts)
* Learn more about available [APIs for managing Logz.io Metrics Accounts](https://api-docs.logz.io/docs/logz/create-a-new-metrics-account)

## Example Usage
```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_metrics_account" "my_metrics_account" {
  email = "user@logz.io"
  account_name = "test"
  plan_uts = 100
  authorized_accounts = [
   12345
  ]
}
```

## Argument Reference

### Required:
* `email` - (String) Email address of an existing admin user on the main account which will also become the admin of the subaccount being created.

### Optional
* `account_name` - (String) Name of the metrics account.
* `plan_uts` - (Integer) Amount of unique time series that can be ingested to the metrics account.
* `authorized_accounts` - (List) IDs of accounts that can access the account's data. Can be an empty array.

##  Attribute Reference
* `account_id` - ID of the metrics account.
* `account_token` - Log shipping token for the metrics account. [Learn more](https://docs.logz.io/user-guide/tokens/log-shipping-tokens/)


## Endpoints used
* [Create](https://api-docs.logz.io/docs/logz/create-a-new-metrics-account).
* [Get](https://api-docs.logz.io/docs/logz/get-a-specific-metrics-account).
* [GetAll](https://api-docs.logz.io/docs/logz/get-a-list-of-all-metrics-accounts).
* [Update](https://api-docs.logz.io/docs/logz/update-a-metrics-account).
* [Delete](https://api-docs.logz.io/docs/logz/delete-a-metrics-account).
