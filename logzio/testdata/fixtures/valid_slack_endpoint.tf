resource "logzio_endpoint" "valid_slack_endpoint" {
  endpoint_type = "slack"
  title = "valid_slack_endpoint"
  description = "some valid description"
  slack {
    url = "https://jsonplaceholder.typicode.com/todos/1"
  }
}

output "test_id" {
  value = logzio_endpoint.valid_slack_endpoint.id
}
