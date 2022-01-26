resource "logzio_endpoint" "%s" {
  title = "my_custom_title"
  endpoint_type = "custom"
  description = "this_is_my_description"
  custom {
    url = "https://jsonplaceholder.typicode.com/todos/1"
    method = "POST"
    headers = "this=is,a=header"
    body_template = jsonencode({
      this: "is"
      my: "template"
    })
  }
}