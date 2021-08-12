resource "logzio_endpoint" "pagerduty" {
  title = "my_pagerduty_title_empty_service_key"
  endpoint_type = "pagerduty"
  description = "this is my description"
  pagerduty {
    service_key = ""
  }
}