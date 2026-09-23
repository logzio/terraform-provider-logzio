# Tracing Account Provider

Provides a Logz.io consumption tracing account resource. This can be used to create and manage
Logz.io tracing accounts on a **Consumption** plan.

* Learn more about available [APIs for managing Logz.io accounts](https://api-docs.logz.io/docs/logz/logz-io-api/)

**Note:** These endpoints are available for **Consumption** accounts only, and are additionally
gated per account by a feature flag. If the account is not eligible the API responds with `400`.

## Example Usage
```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_tracing_account" "my_tracing_account" {
  email         = "user@logz.io"
  account_name  = "test"
  max_daily_gb  = 5
  authorized_accounts = [
    12345
  ]
}
```

## Argument Reference

### Required:
* `email` - (String) Email address of an existing admin user on the main account which will also become the admin of the tracing account being created.
* `account_name` - (String) Name of the tracing account.
* `max_daily_gb` - (Float) The maximum volume of data, in GB, that the tracing account can index per calendar day. This **is** the account's soft cap - the API reads the soft limit from this value and writes it back to this value, so there is no separate soft limit argument. Must not be negative. The account is suspended once its daily usage exceeds this value, so `0` suspends it as soon as it receives any data.

### Optional
* `authorized_accounts` - (List) IDs of accounts that can access the account's data. Can be an empty array.

##  Attribute Reference
* `account_id` - ID of the tracing account.
* `account_token` - Shipping token for the tracing account.
* `retention` - Number of days that trace data is retained.

### Import tracing accounts as resources

```
terraform import logzio_tracing_account.my_tracing_account <TRACING-ACCOUNT-ID>
```
