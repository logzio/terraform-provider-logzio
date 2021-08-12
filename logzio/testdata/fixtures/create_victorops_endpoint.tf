resource "logzio_endpoint" "victorops" {
  title = "my_victorops_title"
  endpoint_type = "victorops"
  description = "this is my description"
  victorops {
    routing_key = "my_routing_key"
    message_type = "my_message_type"
    service_api_key = "my_service_api_key"
  }
}