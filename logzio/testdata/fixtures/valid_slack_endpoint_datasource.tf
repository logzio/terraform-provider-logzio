resource "logzio_endpoint" "valid_slack_endpoint_datasource" {
  endpoint_type = "Slack"
  title = "valid_slack_endpoint_datasource"
  description = "some valid description"
  slack {
    url = "https://www.test.com"
  }
}

data "logzio_endpoint" "by_title" {
  title = "valid_slack_endpoint_datasource"
  depends_on = ["logzio_endpoint.valid_slack_endpoint_datasource"]
}

output "valid_slack_endpoint_datasource_id" {
  value = logzio_endpoint.valid_slack_endpoint_datasource.id
}