# Log Shipping Token Datasource

Use this data source to access information about existing Logz.io log shipping tokens.

* Learn more about log shipping tokens in the [Logz.io Docs](https://docs.logz.io/api/#tag/Manage-log-shipping-tokens).

## Argument Reference

* `name` - (String) Descriptive name for this log shipping token.
* `token_id` - (Integer) The log shipping token's ID.
* `enabled` - (Boolean) To enable this log shipping token, true. To disable, false. **Note:** this argument can only be set after the creation of the token. Each token is created with the `enabled` argument set to true. You can set this field to `false` on update.  

##  Attribute Reference

* `token` - (String) The log shipping token itself.
* `updated_at` - (Integer) Unix timestamp of when this log shipping token was last updated.
* `updated_by` - (String) Email address of the last user to update this log shipping token.
* `created_at` - (Integer) Unix timestamp of when this log shipping token was created.
* `created_by` - (String) Email address of the user who created this log shipping token.
