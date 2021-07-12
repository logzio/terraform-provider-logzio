variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_log_shipping_token" "my_token" {
  name = "hello_world"
}