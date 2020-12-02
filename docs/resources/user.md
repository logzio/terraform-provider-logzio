# User Provider

Provides a Logz.io user resource. This can be used to create and manage Logz.io users.

* Learn more about available [APIs for managing Logz.io users](https://docs.logz.io/api/#tag/Manage-users)

## Example Usage

```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

variable "account_id" {
  description = "the account id you want to use to create your user in"
}

provider "logzio" {
  api_token = var.api_token
}

# Create a new user
resource "logzio_user" "my_user" {
  username = "test.user@this.test"
  fullname = "test user"
  roles = [ 2 ]
  account_id = var.account_id
}
```

## Argument Reference

* `username` - (Optional) Username credential.
* `fullname` - (Optional) First and last name of the user.
* `account_id` - (Optional) Logz.io account ID.
* `roles` - (Required) For User access, `2`. For Admin access, `3`.

##  Attribute Reference

* `id` - (Optional) ID of the user in the Logz.io platform.
* `active` - (Optional) Defaults to `true`. If `true`, the user is active, if `false`, the user is suspended.




## Endpoints used

* [Create user](https://docs.logz.io/api/#operation/createUser)
