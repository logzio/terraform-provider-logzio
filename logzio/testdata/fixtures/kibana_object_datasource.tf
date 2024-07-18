data "logzio_kibana_object" "ds_kb" {
  object_id = "search:tf-provider-datasource-test-search"
  object_type = "search"
  depends_on = ["logzio_kibana_object.test_kb_for_datasource"]
}

output "output_id" {
  value = data.logzio_kibana_object.ds_kb.id
}