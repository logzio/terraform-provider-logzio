resource "logzio_log_shipping_token" "log_shipping_token_datasource" {
  name = "my_token"
}

data "logzio_log_shipping_token" "my_log_shipping_token_datasource" {
  name = "${logzio_log_shipping_token.log_shipping_token_datasource.name}"
  depends_on = ["logzio_log_shipping_token.log_shipping_token_datasource"]
}