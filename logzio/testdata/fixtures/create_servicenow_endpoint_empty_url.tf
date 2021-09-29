resource "logzio_endpoint" "servicenow" {
  title = "servicenow_empty_url"
  endpoint_type = "servicenow"
  description = "this is my description"
  servicenow {
    username = "my_username"
    password = "my_password"
    url = ""
  }
}