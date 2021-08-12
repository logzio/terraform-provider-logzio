resource "logzio_endpoint" "servicenow" {
  title = "servicenow_empty_username"
  endpoint_type = "servicenow"
  description = "this is my description"
  servicenow {
    username = ""
    password = "my_password"
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}