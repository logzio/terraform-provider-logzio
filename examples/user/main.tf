variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

variable "account_id" {
  description = "the account id you want to use to create your user in"
}

provider "logzio" {
  api_token = "${var.api_token}"
}

resource "logzio_user" "my_user" {
  username = "test.user@this.test"
  fullname = "test user"
  roles = [ 2 ]
  account_id = "${var.account_id}"
}