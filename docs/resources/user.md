# User Provider

Provides a Logz.io user resource. This can be used to create and manage Logz.io users.

* Learn more about available [APIs for managing Logz.io users](https://docs.logz.io/api/#tag/Manage-users)

## Example Usage

```hcl
# Create a new user
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


resource "logzio_user" "my_user" {
  username = "test.user@this.test"
  fullname = "test user"
  roles = [ 2 ]
  account_id = var.account_id
}
```

## Argument Reference

* `fullname` - (Required) First and last name of the user.
* `username` - (Required) Username credential.
* `roles` - (Required) For User access, `2`. For Admin access, `3`.
* `active` - (Required) If `true`, the user is active, if `false`, the user is suspended.
* `account_id` - (Required) Logz.io account ID.


##  Attribute Reference

* `id` - ID of the user in the Logz.io platform.




## Endpoints used

* [Create user](https://docs.logz.io/api/#operation/createUser)
