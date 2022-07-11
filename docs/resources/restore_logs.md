# Restore Logs Provider

Provides a Logz.io restore logs resource. This can be used to create and manage Logz.io restore logs operations.

* Learn more about restore in the [Logz.io Docs](https://docs.logz.io/api/#tag/Restore-logs)

**Note:** In order to initiate a restore operation you must have an archive linked to your Logz.io account.
## Example Usage

```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_restore_logs" "my_restore" { 
  account_name = "test_restore"
  start_time = 1635134987
  end_time = 1635145789
}
```

## Argument Reference

* `account_name` - (String) Name of the restored account.
* `start_time` - (Integer) UNIX timestamp in milliseconds specifying the earliest logs to be restored.
* `end_time` - (Integer) UNIX timestamp in milliseconds specifying the latest logs to be restored.

**Note:** Once a restore operation was created, changing any of its arguments will cause the resource to be destroyed
re-created under a new ID.

##  Attribute Reference

* `restore_operation_id` - (Integer) ID of the restore operation in Logz.io.
* `account_id` - (Integer) ID of the restored account in Logz.io.
* `restored_volume_gb` - (Float) Volume of data restored so far. If the restore operation is still in progress, this will be continuously updated.
* `status` - (String) Returns the current status of the restored account. See [documentation](https://docs.logz.io/api/#operation/getRestoreRequestByIdApi) for more info about the possible statuses and their meaning.
* `created_at` - (Integer) Timestamp when the restore process was created and entered the queue. Since only one account can be restored at a time, the process may not initiate immediately.
* `started_at` - (Integer) UNIX timestamp in milliseconds when the restore process initiated.
* `finished_at` - (Integer) UNIX timestamp in milliseconds when the restore process completed.
* `expires_at` - (Integer) UNIX timestamp in milliseconds specifying when the account is due to expire. Restored accounts expire automatically after a number of days, as specified in the account's terms.

## Importing resource:
To import a restore operation you'll need to specify the restore's id.
For example, if you have in your Terraform configuration the following:

```hcl
resource "logzio_restore_logs" "imported" {
  ...
}
```

And your restore operation id is 123456, then your import command should be:

```bash
terraform import logzio_restore_logs.imported 123456
```