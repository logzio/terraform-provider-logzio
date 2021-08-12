# Subaccount Provider

Provides a Logz.io Log Management subaccount resource. This can be used to create and manage Logz.io log monitoring subaccounts.

* Learn more about accounts in the [Logz.io Docs](https://docs.logz.io/user-guide/accounts/manage-the-main-account-and-sub-accounts.html)
* Learn more about available [APIs for managing Logz.io subaccounts](https://docs.logz.io/api/#tag/Manage-sub-accounts)

## Example Usage

```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_subaccount" "my_subaccount" {
  email = "user@logz.io"
  account_name = "test"
  retention_days = 2
  sharing_objects_accounts = [
    12345
  ]
  frequency_minutes = 3
  utilization_enabled = true
}
```

## Argument Reference

### Required:
* `email` - (String) Email address of an existing admin user on the main account which will also become the admin of the subaccount being created.
* `account_name` - (String) Name of the subaccount.
* `max_daily_gb` - (Float) Maximum daily log volume that the subaccount can index, in GB.
* `retention_days` - (Integer) Number of days that log data is retained.
* `sharing_objects_accounts` - (List) IDs of accounts that can access the account's Kibana objects. Can be an empty array.

### Optional
* `searchable` - (Boolean) False by default. Determines if other accounts can search logs indexed by the subaccount.
* `accessible` - (Boolean) False by default. Determines if users of main account can access the subaccount.
* `doc_size_setting` - (Boolean) False by default. If enabled, Logz.io adds a `LogSize` field to record the size of the log line in bytes, for the purpose of managing account utilization. [Learn more about managing account usage](https://docs.logz.io/user-guide/accounts/manage-account-usage.html#enabling-account-utilization-metrics-and-log-size)
* `utilization_enabled` - (Boolean) If enabled, account utilization metrics and expected utilization at the current indexing rate are measured at a pre-defined sampling rate. Useful for managing account utilization and avoiding running out of capacity. [Learn more about managing account capacity](https://docs.logz.io/user-guide/accounts/manage-account-usage.html).
* `frequency_minutes` - (Int) Determines the sampling rate in minutes of the utilization.
* `flexible` - (Boolean) Defaults to false. Whether the sub account that created is flexible or not. Can be set to flexible only if the main account is flexible.
* `reserved_daily_gb` - (Float) The maximum volume of data that an account can index per calendar day. Depends on `flexible`. For further info see [the docs](https://docs.logz.io/api/#operation/createTimeBasedAccount).

##  Attribute Reference

* `account_id` - ID of the subaccount.
* `account_token` - Log shipping token for the subaccount. [Learn more](https://docs.logz.io/user-guide/tokens/log-shipping-tokens/)

**Note:** The above attributes displayed only from v1.2.4. If you're using an earlier version, please upgrade and use `terraform apply -refersh` to add those attributes to your existing resources.

## Endpoints used
* [Create](https://docs.logz.io/api/#operation/createTimeBasedAccount).
* [Get](https://docs.logz.io/api/#operation/get).
* [GetAll](https://docs.logz.io/api/#operation/getAll).
* [Update](https://docs.logz.io/api/#operation/updateTimeBasedAccount).
* [Delete](https://docs.logz.io/api/#operation/deleteTimeBasedAccount).