resource "logzio_endpoint" "opsgenie" {
  title = "my_opsgenie_title"
  endpoint_type = "opsgenie"
  description = "this is my description"
  opsgenie {
    api_key = "my_api_key"
  }
}