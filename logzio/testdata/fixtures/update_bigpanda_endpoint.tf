resource "logzio_endpoint" "bigpanda" {
  title = "bigpanda_title_updated"
  endpoint_type = "bigpanda"
  description = "this is my description"
  bigpanda {
    api_token = "updated_api_token"
    app_key = "updated_app_key"
  }
}