variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_kibana_object" "my_search" {
  kibana_version = "7.2.1"
  data = file("./search.json")
}