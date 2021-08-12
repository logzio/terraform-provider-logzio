resource "logzio_endpoint" "servicenow" {
  title = "my_servicenow_title"
  endpoint_type = "servicenow"
  description = "this is my description"
  servicenow {
    username = "my_username"
    password = "my_password"
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}