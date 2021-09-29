resource "logzio_endpoint" "datadog" {
  title = "my_datadog_title"
  endpoint_type = "datadog"
  description = "this is my description"
  datadog {
    api_key = "my_api_key"
  }
}