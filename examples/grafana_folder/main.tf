variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource logzio_grafana_folder "my_folder" {
  title = "my_title"
}