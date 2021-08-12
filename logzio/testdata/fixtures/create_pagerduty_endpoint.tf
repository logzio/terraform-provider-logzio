resource "logzio_endpoint" "pagerduty" {
  title = "my_pagerduty_title"
  endpoint_type = "pagerduty"
  description = "this is my description"
  pagerduty {
    service_key = "my_service_key"
  }
}