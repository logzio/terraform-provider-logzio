variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_grafana_dashboard" "my_dashboard" {
  dashboard_json = <<EOD
{
  "title": "a title",
  "uid": "my_dashboard_uid",
  "panels": []
}
EOD
  folder_uid = "my_folder_uid"
  message = "my message"
  overwrite = true
}