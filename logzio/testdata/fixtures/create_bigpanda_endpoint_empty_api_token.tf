resource "logzio_endpoint" "bigpanda" {
  title = "big_panda_empty_api_token"
  endpoint_type = "bigpanda"
  description = "this is my description"
  bigpanda {
    api_token = ""
    app_key = "my_app_key"
  }
}