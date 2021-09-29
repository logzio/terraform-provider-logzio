resource "logzio_endpoint" "servicenow" {
  title = "updated_servicenow_title"
  endpoint_type = "servicenow"
  description = "this is my description"
  servicenow {
    username = "updated_username"
    password = "updated_password"
    url = "https://jsonplaceholder.typicode.com/todos/2"
  }
}