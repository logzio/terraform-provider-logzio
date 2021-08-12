resource "logzio_endpoint" "victorops" {
  title = "victorops_empty_routing_key"
  endpoint_type = "victorops"
  description = "this is my description"
  victorops {
    routing_key = ""
    message_type = "my_message_type"
    service_api_key = "my_service_api_key"
  }
}