resource "logzio_endpoint" "bigpanda" {
  title = "big_panda_empty_app_key"
  endpoint_type = "bigpanda"
  description = "this is my description"
  bigpanda {
    api_token = "my_api_token"
    app_key = ""
  }
}