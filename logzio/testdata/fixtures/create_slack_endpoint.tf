resource "logzio_endpoint" "%s" {
  endpoint_type = "Slack"
  title = "slack_endpoint"
  description = "some valid description"
  slack {
    url = "https://www.test.com"
  }
}