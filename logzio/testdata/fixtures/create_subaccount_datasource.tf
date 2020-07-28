resource "logzio_subaccount" "subaccount_datasource" {
  email = "your@email.com"
  account_name = "test"
  retention_days = 2
  max_daily_gb = 1
  sharing_objects_accounts = [
    12345
  ]
}

data "logzio_subaccount" "subaccount_datasource_by_id" {
  account_id = logzio_subaccount.subaccount_datasource.id
  depends_on = ["logzio_subaccount.subaccount_datasource"]
}