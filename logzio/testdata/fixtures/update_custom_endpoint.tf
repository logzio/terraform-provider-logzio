resource "logzio_endpoint" "%s" {
  title = "updated_custom_endpoint"
  endpoint_type = "custom"
  description = "this_is_my_description"
  custom {
    url = "https://jsonplaceholder.typicode.com/todos/2"
    method = "PUT"
    body_template = jsonencode({})
    headers = ""
  }
}