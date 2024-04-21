variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_metrics_account" "my_metrics_account" {
  email = "user@logz.io"
  account_name = "test"
  plan_uts = 100
  authorized_accounts = [
   12345
  ]
}
