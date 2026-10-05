# Warm Tier Provider

Provides a Logz.io warm tier resource. This can be used to manage the warm tier retention of an existing Logz.io log account: the main account of the API token, or one of its sub accounts.
It doesn't create or delete accounts.

* Learn more about available [APIs for managing Logz.io accounts](https://api-docs.logz.io/docs/logz/logz-io-api/)

**Note:** The main account must already have the warm tier enabled. Turning the warm tier of the main account on or off is done by Logz.io support.

## Example Usage
```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_warm_tier" "main_account" {
  account_id                 = 12345
  snap_search_retention_days = 7
}

resource "logzio_subaccount" "my_subaccount" {
  email          = "user@logz.io"
  account_name   = "test"
  retention_days = 4
  max_daily_gb   = 1
  sharing_objects_accounts = [
    12345
  ]
}

resource "logzio_warm_tier" "my_subaccount" {
  account_id                 = logzio_subaccount.my_subaccount.account_id
  snap_search_retention_days = 3
}
```

## Argument Reference

### Required:
* `account_id` - (Int) ID of the main account of the API token, or of one of its sub accounts. Changing it creates a new resource.
* `snap_search_retention_days` - (Int) Number of days to retain data in the warm tier. For the main account, at least `1`. For a sub account, from `0` (no warm tier) up to the main account's warm tier retention. Lowering the value deletes the warm tier data older than the new retention.

## Usage with `logzio_subaccount`

To manage the warm tier of a sub account with this resource, leave `snap_search_retention_days` out of its `logzio_subaccount` resource, and reference the sub account's `account_id` as in the example above. The reference makes Terraform apply the sub account first.
A `logzio_subaccount` without `snap_search_retention_days` keeps the warm tier retention the account has when it's updated.

## Destroy

Destroying the resource only removes it from the Terraform state. The account keeps its warm tier retention.

### Import warm tier as a resource

```
terraform import logzio_warm_tier.my_warm_tier <ACCOUNT-ID>
```
