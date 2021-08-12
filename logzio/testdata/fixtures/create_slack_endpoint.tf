resource "logzio_endpoint" "%s" {
  endpoint_type = "slack"
  title = "slack_endpoint"
  description = "some valid description"
  slack {
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}