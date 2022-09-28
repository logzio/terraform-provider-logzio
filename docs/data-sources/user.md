# User Datasource

Use this data source to access information about existing Logz.io users.

* Learn more about available [APIs for managing Logz.io users](https://docs.logz.io/api/#tag/Manage-users)

## Argument Reference

* `id` - ID of the user in the Logz.io platform.
* `username` - Username credential.

##  Attribute Reference

* `fullname` - First and last name of the user.
* `role` - User role. Can be `USER_ROLE_READONLY`, `USER_ROLE_REGULAR` or `USER_ROLE_ACCOUNT_ADMIN`.
* `active` - If `true`, the user is active, if `false`, the user is suspended.
* `account_id` - Logz.io account ID.