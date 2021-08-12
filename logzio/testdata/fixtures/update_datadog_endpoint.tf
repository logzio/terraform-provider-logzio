resource "logzio_endpoint" "datadog" {
  title = "datadog_title_updated"
  endpoint_type = "datadog"
  description = "this is my description"
  datadog {
    api_key = "updated_api_key"
  }
}