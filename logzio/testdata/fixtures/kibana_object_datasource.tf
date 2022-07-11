resource "logzio_kibana_object" "test_kb_for_datasource" {
  kibana_version = "7.2.1"
  data = file("./testdata/fixtures/kibana_objects/create_search.json")
}

data "logzio_kibana_object" "ds_kb" {
  object_id = "search:tf-provider-test-search"
  object_type = "search"
  depends_on = ["logzio_kibana_object.test_kb_for_datasource"]
}

output "output_id" {
  value = data.logzio_kibana_object.ds_kb.id
}