resource "logzio_endpoint" "bigpanda" {
  title = "my_bigpanda_title"
  endpoint_type = "bigpanda"
  description = "this is my description"
  bigpanda {
    api_token = "my_api_token"
    app_key = "my_app_key"
  }
}