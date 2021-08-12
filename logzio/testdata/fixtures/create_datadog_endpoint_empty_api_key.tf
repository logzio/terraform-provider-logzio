resource "logzio_endpoint" "datadog" {
  title = "datadog_empty_api_key"
  endpoint_type = "datadog"
  description = "this is my description"
  datadog {
    api_key = ""
  }
}