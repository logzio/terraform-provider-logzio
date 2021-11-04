resource "logzio_authentication_groups" "my_auth_groups" {
  authentication_group {
    group = "group_testing_1"
    user_role = "USER_ROLE_ACCOUNT_ADMIN"
  }
  authentication_group {
    group = "group_testing_2"
    user_role = "USER_ROLE_REGULAR"
  }
  authentication_group {
    group = "group_testing_3"
    user_role = "USER_ROLE_READONLY"
  }
}
