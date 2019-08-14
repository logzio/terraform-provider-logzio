resource "logzio_subaccount" "create_subaccount" {
  email = "jon.boydell@massive.co"
  account_name = "create_subaccount_name"
  max_daily_gb = "10.0"
  retention_days = "10"
  accessible = true
  searchable = true
  sharing_objects_account = []
  doc_size_setting = true
}