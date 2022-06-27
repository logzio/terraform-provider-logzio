resource "logzio_kibana_object" "test_kb_obj_search" {
  kibana_version = "7.2.1"
  data = file("./create_search.json")
}