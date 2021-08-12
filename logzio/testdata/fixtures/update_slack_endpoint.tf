resource "logzio_endpoint" "%s" {
  endpoint_type = "slack"
  title = "updated_slack_endpoint"
  description = "some updated description"
  slack {
    url = "https://jsonplaceholder.typicode.com/todos/2"
  }
}
