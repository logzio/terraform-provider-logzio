resource "logzio_kibana_object" "test_kb_obj_dashboard" {
  kibana_version = "7.2.1"
  data = file("./testdata/fixtures/kibana_objects/create_dashboard.json")
}