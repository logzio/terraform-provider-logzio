resource "logzio_endpoint" "pagerduty" {
  title = "pagerduty_title_updated"
  endpoint_type = "pagerduty"
  description = "this is my description"
  pagerduty {
    service_key = "another_service_key"
  }
}