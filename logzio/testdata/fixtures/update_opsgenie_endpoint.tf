resource "logzio_endpoint" "opsgenie" {
  title = "updated_opsgenie_title"
  endpoint_type = "opsgenie"
  description = "this is my description"
  opsgenie {
    api_key = "updated_api_key"
  }
}