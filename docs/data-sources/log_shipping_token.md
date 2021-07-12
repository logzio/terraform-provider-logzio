# Log Shipping Token Provider

Provides a Logz.io log shipping token resource. This can be used to create and manage Logz.io log shipping tokens.

* Learn more about log shipping tokens in the [Logz.io Docs](https://docs.logz.io/api/#tag/Manage-log-shipping-tokens)

## Example Usage

```hcl
variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_log_shipping_token" "my_token" {
  name = "hello_world"
}
```

## Argument Reference

* `name` - (String) Descriptive name for this token.
* `enabled` - (Boolean) To enable this token, true. To disable, false. **Note:** this argument can only be set after the creation of the token. Each token is created with the `enabled` argument set to true. You can set this field to `false` on update.  

##  Attribute Reference

* `token_id` - (Integer) The token's ID.
* `token` - (String) The token itself.
* `updated_at` - (Integer) Unix timestamp of when this token was last updated.
* `updated_by` - (String) Email address of the last user to update this token.
* `created_at` - (Integer) Unix timestamp of when this token was created.
* `created_by` - (String) Email address of the user who created this token.