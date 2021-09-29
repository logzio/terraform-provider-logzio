resource "logzio_endpoint" "opsgenie" {
  title = "opsgenie_empty_api_key"
  endpoint_type = "opsgenie"
  description = "this is my description"
  opsgenie {
    api_key = ""
  }
}