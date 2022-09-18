variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_user" "my_user" {
  username = "test.user@this.test"
  fullname = "test user"
  role = "USER_ROLE_READONLY"
  account_id = 1234
}