resource "logzio_kibana_object" "test_kb_obj_search" {
  kibana_version = "7.2.1"
  data = file("./testdata/fixtures/kibana_objects/create_search.json")
}