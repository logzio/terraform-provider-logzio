resource "logzio_endpoint" "valid_slack_endpoint" {
  endpoint_type = "Slack"
  title = "valid_slack_endpoint"
  description = "some valid description"
  slack {
    url = "https://www.test.com"
  }
}

output "test_id" {
  value = logzio_endpoint.valid_slack_endpoint.id
}