resource "logzio_endpoint" "invalid_slack_endpoint" {
  endpoint_type = "slack"
  title = "invalid_slack_endpoint"
  description = "some description"
  slack {
    url = "some/bad/url"
  }
}
