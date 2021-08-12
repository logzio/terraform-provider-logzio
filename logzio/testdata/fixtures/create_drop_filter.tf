resource "logzio_drop_filter" "%s" {
  log_type = "some_type_create"

  field_conditions {
    field_name = "some_field"
    value = "some_string_value"
  }
  field_conditions {
    field_name = "another_field"
    value = 200
  }
}
