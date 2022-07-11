resource "logzio_kibana_object" "test_kb_obj_visualization" {
  kibana_version = "7.2.1"
  data = file("./testdata/fixtures/kibana_objects/update_visualization.json")
}