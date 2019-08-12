resource "logzio_endpoint" "%s" {
  endpoint_type = "Slack"
  title = "updated_slack_endpoint"
  description = "some updated description"
  slack {
    url = "https://www.test.com"
  }
}