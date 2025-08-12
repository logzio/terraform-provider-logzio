# Metrics Drop Filter Provider

Provides a Logz.io metrics drop filter datasource. This can be used to get and manage Logz.io metrics drop filters.

* Learn more about drop filters in the [Logz.io Docs](https://docs.logz.io/docs/user-guide/data-hub/drop-filters/drop-fiters-metrics/).

## Argument Reference

* `account_id` - The Logz.io metrics account ID to which the drop filter applies.
* `active` - If true, the drop filter is active and metrics that match the filter are dropped before indexing. If false, the drop filter is disabled. Defaults to true.
* `filter` - The filter object that defines the drop filter criteria. See below for nested schema.
    * `name` - The name of the metric label to filter on.
    * `value` - The value of the metric label to match against.
    * `condition` - The comparison operator to use for the filter. Supported values are `EQ`, `NOT_EQ`, `REGEX_MATCH`, and `REGEX_NO_MATCH`.
* `operator` - The logical operator for combining filter expressions. Supported value is `AND`.

## Attribute Reference

* `id` - (String) The unique identifier of the drop filter in the Logz.io database.
* `created_at` - (String) The timestamp when the drop filter was created.
* `updated_at` - (String) The timestamp when the drop filter was last updated.
* `created_by` - (String) The user who created the drop filter.
* `modified_at` - (String) The timestamp when the drop filter was last modified.
* `modified_by` - (String) The user who last updated the drop filter.
