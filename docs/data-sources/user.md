# User Resource

Use this data source to access information about existing Logz.io users.

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

* `roles` - (Required) For User access, `2`. For Admin access, `3`.
* `active` - (Required) If `true`, the user is active, if `false`, the user is suspended.
* `account_id` - (Optional) Logz.io account ID.
* `fullname` - (Optional) First and last name of the user.


##  Attribute Reference

* `id` - (Optional) ID of the user in the Logz.io platform.
* `username` - (Optional) Username credential.

## Endpoints used

* [Get all users in main account and associated subaccounts](https://docs.logz.io/api/#operation/listAllAccountUsers)
* [Get all users in account](https://docs.logz.io/api/#operation/listUsers)
* [Update user](https://docs.logz.io/api/#operation/updateUser)
* [Delete user from account](https://docs.logz.io/api/#operation/deleteUser)
* [Delete user from all accounts](https://docs.logz.io/api/#operation/deleteUserRecursively)
* [Get user by ID](https://docs.logz.io/api/#operation/getUser)

