resource "logzio_endpoint" "victorops" {
  title = "updated_victorops_title"
  endpoint_type = "victorops"
  description = "this is my description"
  victorops {
    routing_key = "updated_routing_key"
    message_type = "updated_message_type"
    service_api_key = "updated_service_api_key"
  }
}