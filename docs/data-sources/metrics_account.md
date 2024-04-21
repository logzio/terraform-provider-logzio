# Subaccount Datasource

Use this data source to access information about existing Logz.io Metrics Accounts.

* Learn more about accounts in the [Logz.io Docs](https://docs.logz.io/docs/user-guide/admin/logzio-accounts/manage-the-main-account-and-sub-accounts).
* Learn more about available [APIs for managing Logz.io Metrics accounts](https://api-docs.logz.io/docs/logz/create-a-new-metrics-account).

## Argument Reference

* `account_id` - ID of the metrics account.

##  Attribute Reference

* `email` - (String) Email address of an existing admin user on the main account which will also become the admin of the created metrics account.
* `account_name` - (String) Name of the metrics account.
* `plan_uts` - (Integer) Amount of unique time series that can be ingested to the metrics account.
* `authorized_accounts` - (List) IDs of accounts that can access the account's data. Can be an empty array.
