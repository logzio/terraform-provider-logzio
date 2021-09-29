resource "logzio_endpoint" "victorops" {
  title = "victorops_empty_service_key"
  endpoint_type = "victorops"
  description = "this is my description"
  victorops {
    routing_key = "my_routing_key"
    message_type = "my_message_type"
    service_api_key = ""
  }
}