resource "logzio_endpoint" "servicenow" {
  title = "servicenow_empty_password"
  endpoint_type = "servicenow"
  description = "this is my description"
  servicenow {
    username = "my_username"
    password = ""
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}