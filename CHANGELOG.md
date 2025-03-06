# Changes by Version

<!-- next version -->
## v1.17.1
- Add missing platform builds for `windows/arm64` and `darwin/arm64`, which were not included in `v1.17.0`

## v1.17.0
- Upgrade `logzio_client_terraform` to `1.23.2`.
- Upgrade Grafana API to match Grafana 10.
    - **Breaking changes:**
        - Grafana Alerts API `exec_err_state` field is no longer configurable, always defaults to `OK`.
        - Grafana Contact Point `teams` now refers to Microsoft Teams Workflows Contact Point. The old `teams` endpoint was deprecated.
        - Default notification policy contact point changed from `default-email` to `grafana-default-email`
    - Webhook Contact Point now supports `authorization_credentials` 
- Upgrade to Go 1.24.

## v1.16.0
- Upgrade `logzio_client_terraform` to `1.22.0`.
- Validate account existence before updating and set an empty account name If the name did not change to prevent API errors

## v1.15.0
- Upgrade `logzio_client_terraform` to `1.21.0`.
- Support [Metrics Accounts API](https://api-docs.logz.io/docs/logz/create-a-new-metrics-account)

## v1.14.1
- Upgrade `logzio_client_terraform` to `1.20.1`.

## v1.14.0
- Upgrade `logzio_client_terraform` to `1.20.0`.
- Support Grafana folders.
- Support [Grafana Alert Rules](https://api-docs.logz.io/docs/logz/get-alert-rules).
- Support [Grafana Contact Point](https://api-docs.logz.io/docs/logz/route-get-contactpoints).
- Grafana Notification Policy.

## v1.13.2
- **alerts v2**:
    - **Bug fix**: Add default value of `JSON` to `output_type` to match API behaviour.
- **grafana dashboards**:
    - **Bug fix**: Handle dashboard configs with id and version fields.

## v1.13.1
- **s3 fetcher**:
    - **Bug fix**: `prefix` field now being applied.

## v1.13.0
- Upgrade `logzio_client_terraform` to `1.16.0`.
- Support [Grafana Dashboards API](https://api-docs.logz.io/docs/logz/get-datasource-by-account).

## v1.12.0
- Upgrade `logzio_client_terraform` to `1.15.0`.
- Support [S3 Fetcher API](https://api-docs.logz.io/docs/logz/create-buckets).

## v1.11.0
- Upgrade `terraform-plugin-sdk` to `v2.24.1`.
- Upgrade `logzio_client_terraform` to `1.14.0`.
- `alert_v2` - support alert scheduling.

## v1.10.2
- Upgrade to logzio_client_terraform 1.13.1.
- Remove retries resources, as the new client handles it. Retries on update that validate update are still active.

## v1.10.1
- Upgrade to logzio_client_terraform 1.13.0.
- Bug fix for **subaccount** - field `reserved_daily_gb` can be 0.

## v1.10.0
- **Breaking Changes**:
    - Upgrading `terraform-plugin-sdk` to `v2.21.0`:
        - **Terraform 0.11** and lower will not be supported.
        - To read more about migrating to v2 of the `terraform-plugin-sdk` see [this article](https://www.terraform.io/plugin/sdkv2/guides/v2-upgrade-guide).
        - Removal of the `logzio_alert` resource. Use `logzio_alert_v2` instead.
        - `logzio_user`:
            - Update resource `logzio_user` to match current API and repo code conventions
            - Removal of `roles` field. Use field `role` instead.
        - `logzio_archive_logs`:
            - Removal of `s3_secret_credentials`. Use `aws_access_key` and `aws_secret_key` instead. See documentation for current configuration structure.
            - Removal of `archiveLogsAzureBlobStorageSettings`, `archiveLogsAmazonS3StorageSettings`. Fields under these attributes are not nested anymore. See docs & examples for reference.
            - Renaming fields.
            - Refactor code to match current API.
            - datasource - removal of secret fields. See documentation for current available fields.
        - `logzio_restore_logs`:
            - Add `username` field, to match current API.
        - `logzio_subaccount`:
            - Removal of field `utilization_settings`.
        - Delete resource from state on unsuccessful read operation.
- Upgrade to Go 1.18.
- Upgrade to logzio_client_terraform 1.12.0.

## v1.9.2
- *Bug fix*: Fix diff for resource `alert_v2` in fields `alert_notification_endpoints`, `notification_emails` ([#116](https://github.com/logzio/terraform-provider-logzio/issues/116)).

## v1.9.1
- *Bug fix*: plugin won't crash when import for `archive_logs` fails.

## v1.9.0
- Update client version(v1.11.0).
- Support [Kibana objects](https://api-docs.logz.io/docs/logz/import-or-export-kibana-objects)

## v1.8.3
- Update client version(v1.10.3).
- Bug fixes:
    - **alerts_v2**:
        - Fix noisy diff for tags.
        - Field `is_enabled` defaults to `true`.
- **sub_accounts**: allow creating flexible account without `max_daily_gb`.

## v1.8.2
- Update client version(v1.10.2).
- Bug fixes:
    - **alerts_v2**: fix bug for columns requiring sort field.
    - **sub_accounts**: add backoff for creating and updating sub accounts.

## v1.8.1
- Upgrade provider's Go version from 1.15 to 1.16 in code and in release workflow.
- Improve tests - add sleep after each test.

## v1.8.0
- **Breaking change**: **custom endpoint** - refactor Headers - now a string of comma-seperated key-value pairs.
- Update client version (v1.10.1) - bug fix for empty Header field.
- Add to custom endpoint datasource Description field.

## v1.7.0
- Update client version (v1.10).
- Support [authentication groups resource](https://api-docs.logz.io/docs/logz/authentication-groups).
- `alerts_v2`: fix noisy diff for `severity_threshold_tiers`.


## v1.6.1
- Update client version (v1.9.1) - bug fix for not found messages.

## v1.6
- Update client version (v1.9).
- Support [archive logs resource](https://api-docs.logz.io/docs/logz/archive-logs).
- Support [restore logs resource](https://api-docs.logz.io/docs/logz/restore-logs).

## v1.5
- Update client version(v1.8).
- `sub_account`:
    - **Breaking changes:**
        - Deprecated attribute `utilization_settings`. Use `frequency_minutes` and `utilization_enabled` instead.
        - Added attributes `flexible` and `reserved_daily_gb`.
        - Refactor tests.
        - Refactor code to use Terraform's retry.
- `endpoint`:
    - **Breaking changes:**
        - Types naming.
    - New endpoints for OpsGenie, ServiceNow and Microsoft Teams.
    - Fix bug for `body_template` for `custom` type ([#70](https://github.com/logzio/terraform-provider-logzio/issues/70)).
    - Refactor tests.
    - Refactor code to use Terraform's retry.
- `alerts_v2`:
    - Fix bug for `filter_must`,`filter_must_not` ([#82](https://github.com/logzio/terraform-provider-logzio/issues/82))
    - Refactor tests.
    - Refactor code to use Terraform's retry.
- `drop_filter`:
    - Improve tests.

## v1.4
- Update client version(v1.7).
- Support [Drop Filter](https://api-docs.logz.io/docs/logz/drop-filters) resource.

## v1.3
- Update client version(v1.6).
- Support Log Shipping Token resource.

## v1.2.4
- Update client version(v1.5.3).
- Fix `sub account` to return attributes `account_token` & `account_id`.

## v1.2.3
- Fix bug for `custom endpoint` empty headers.
- Allow empty sharing accounts array in `sub account`.
- Add retry in resource `sub account`.
- Replace module `terraform` with `terraform-plugin-sdk`. See further explanation [here](https://www.terraform.io/docs/extend/guides/v1-upgrade-guide.html).
- Upgrade to Go v1.15.
- Update client version(v1.5.2).

## v1.2.2
- Update client version(v1.5.1).
- Fix alerts_v2 sort bug.

## v1.2.1
- Fix alerts_v2 type "TABLE" bug.

## v1.2
- Update client version(v1.5.0).
- Support Alerts v2 resource.
- Fix 404 error for Alerts.

## v1.1.8
- Update client version
- Fix custom endpoint headers bug

## v1.1.7
- Published to Terraform registry

## v1.1.5
- Fix boolean parameters not parsed bug
- Support import command to state

## v1.1.4
- Support Sub Accounts resource
- few bug fixes
- removed circleCI

## v1.1.3
- examples now use TF12
- will now generate the meta data needed for the IntelliJ type IDE HCL plugin
- no more travis - just circle CI
- version bump to use the latest TF library (0.12.6), now compatible with TF12

## v1.1.2
- Moved some of the source code around to comply with TF provider layout convention
- Moved the examples into an examples directory
