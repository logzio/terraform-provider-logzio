variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_subaccount" "my_subaccount" {
  email = "user@logz.io"
  account_name = "test"
  retention_days = 2
  sharing_objects_accounts = [
    12345
  ]
  frequency_minutes = 3
  utilization_enabled = true
}
