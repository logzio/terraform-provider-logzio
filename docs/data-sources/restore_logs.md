# Restore Logs Datasource

Use this data source to access information about existing Logz.io restore operations.

* Learn more about restore logs in the [Logz.io Docs](https://docs.logz.io/api/#tag/Restore-logs)

## Argument Reference

* `restore_operation_id` - (Integer) ID of the restore operation in Logz.io.

## Attribute Reference

* `account_name` - (String) Name of the restored account.
* `start_time` - (Integer) UNIX timestamp in milliseconds specifying the earliest logs to be restored.
* `end_time` - (Integer) UNIX timestamp in milliseconds specifying the latest logs to be restored.
* `account_id` - (Integer) ID of the restored account in Logz.io.
* `restored_volume_gb` - (Float) Volume of data restored so far. If the restore operation is still in progress, this will be continuously updated.
* `status` - (String) Returns the current status of the restored account. See [documentation](https://docs.logz.io/api/#operation/getRestoreRequestByIdApi) for more info about the possible statuses and their meaning.
* `created_at` - (Integer) Timestamp when the restore process was created and entered the queue. Since only one account can be restored at a time, the process may not initiate immediately.
* `started_at` - (Integer) UNIX timestamp in milliseconds when the restore process initiated.
* `finished_at` - (Integer) UNIX timestamp in milliseconds when the restore process completed.
* `expires_at` - (Integer) UNIX timestamp in milliseconds specifying when the account is due to expire. Restored accounts expire automatically after a number of days, as specified in the account's terms.
