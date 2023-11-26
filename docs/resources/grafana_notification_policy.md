# Grafana Notification Policy Provider

Provides a Logz.io Grafana Notification Policy resource. This can be used to create and manage Grafana notification policies in Logz.io.

### Important Note:

Please note that due to the API limitations, ONE resource of Grafana Notification Policy manages the entire policy tree.
Deleting the resource will reset your ENTIRE notification policy tree.

## Example Usage

```hcl
resource logzio_grafana_notification_policy test_np {
  contact_point = "default-email"
  group_by = ["p8s_logz_name"]
  group_wait      = "50s"
  group_interval  = "7m"
  repeat_interval = "4h"

  policy {
    matcher {
      label = "some_label"
      match = "="
      value = "some_value"
    }
    contact_point = "default-email"
    continue      = true

    group_wait      = "50s"
    group_interval  = "7m"
    repeat_interval = "4h"
    mute_timings = ["some-mute-timing"]

    policy {
      matcher {
        label = "another_label"
        match = "="
        value = "another_value"
      }
      contact_point = "default-email"
    }
  }
}
```

## Argument Reference

### Required:

* `contact_point` - (String) The default contact point to route all unmatched notifications to.
* `group_by` - (List of strings) A list of alert labels to group alerts into notifications by.

### Optional:

* `group_interval` - (String) Minimum time interval between two notifications for the same group.
* `group_wait` - (String) Time to wait to buffer alerts of the same group before sending a notification.
* `repeat_interval` - (String) Minimum time interval for re-sending a notification if an alert is still firing.
* `policy` - (Block List) Routing rules for specific label sets. See below for **nested schema**.

### Nested schema for `policy`:

#### Required:

* `contact_point` - (String) The contact point to route notifications that match this rule to.

#### Optional:

* `group_by` - (List of String) A list of alert labels to group alerts into notifications by.
* `mute_timings` - (List of String) A list of mute timing names to apply to alerts that match this policy. **Warning** - insert names of mute-timing that already exists, otherwise it can cause problems in your system.
* `continue` - (Boolean) Whether to continue matching subsequent rules if an alert matches the current rule. Otherwise, the rule will be 'consumed' by the first policy to match it.
* `group_wait` - (String) Time to wait to buffer alerts of the same group before sending a notification.
* `group_interval` - (String) Minimum time interval between two notifications for the same group.
* `repeat_interval` - (String) Minimum time interval for re-sending a notification if an alert is still firing.
* `matcher` - (Block List) Describes which labels this rule should match. When multiple matchers are supplied, an alert must match ALL matchers to be accepted by this policy. When no matchers are supplied, the rule will match all alert instances. See below for **nested schema**.
* `policy` - (Block List) Routing rules for specific label sets. See below for **nested schema**.

### Nested schema for `policy.matcher`:

#### Required:

* `label` - (String) The name of the label to match against.
* `match` - (String) The operator to apply when matching values of the given label. Allowed operators are `=` (for equality), `!=` (for negated equality), `=~` (for regex equality), and `!~` (for negated regex equality).
* `value` -  (String) The label value to match against.

### Nested schema for `policy.policy`:

#### Required:

* `contact_point` - (String) The contact point to route notifications that match this rule to.

#### Optional:

* `group_by` - (List of String) A list of alert labels to group alerts into notifications by.
* `mute_timings` - (List of String) A list of mute timing names to apply to alerts that match this policy. **Warning** - insert names of mute-timing that already exists, otherwise it can cause problems in your system.
* `continue` - (Boolean) Whether to continue matching subsequent rules if an alert matches the current rule. Otherwise, the rule will be 'consumed' by the first policy to match it.
* `group_wait` - (String) Time to wait to buffer alerts of the same group before sending a notification.
* `group_interval` - (String) Minimum time interval between two notifications for the same group.
* `repeat_interval` - (String) Minimum time interval for re-sending a notification if an alert is still firing.
* `matcher` - (Block List) Describes which labels this rule should match. When multiple matchers are supplied, an alert must match ALL matchers to be accepted by this policy. When no matchers are supplied, the rule will match all alert instances. See below for **nested schema**.
* `policy` - (Block List) Routing rules for specific label sets. See below for **nested schema**.

### Nested schema for `policy.policy.matcher`:

#### Required:

* `label` - (String) The name of the label to match against.
* `match` - (String) The operator to apply when matching values of the given label. Allowed operators are `=` (for equality), `!=` (for negated equality), `=~` (for regex equality), and `!~` (for negated regex equality).
* `value` -  (String) The label value to match against.

### Nested schema for `policy.policy.policy`:

#### Required:

* `contact_point` - (String) The contact point to route notifications that match this rule to.

#### Optional:

* `group_by` - (List of String) A list of alert labels to group alerts into notifications by.
* `mute_timings` - (List of String) A list of mute timing names to apply to alerts that match this policy. **Warning** - insert names of mute-timing that already exists, otherwise it can cause problems in your system.
* `continue` - (Boolean) Whether to continue matching subsequent rules if an alert matches the current rule. Otherwise, the rule will be 'consumed' by the first policy to match it.
* `group_wait` - (String) Time to wait to buffer alerts of the same group before sending a notification.
* `group_interval` - (String) Minimum time interval between two notifications for the same group.
* `repeat_interval` - (String) Minimum time interval for re-sending a notification if an alert is still firing.
* `matcher` - (Block List) Describes which labels this rule should match. When multiple matchers are supplied, an alert must match ALL matchers to be accepted by this policy. When no matchers are supplied, the rule will match all alert instances. See below for **nested schema**.
* `policy` - (Block List) Routing rules for specific label sets. See below for **nested schema**.

### Nested schema for `policy.policy.policy.matcher`:

#### Required:

* `label` - (String) The name of the label to match against.
* `match` - (String) The operator to apply when matching values of the given label. Allowed operators are `=` (for equality), `!=` (for negated equality), `=~` (for regex equality), and `!~` (for negated regex equality).
* `value` -  (String) The label value to match against.

### Nested schema for `policy.policy.policy.policy`:

#### Required:

* `contact_point` - (String) The contact point to route notifications that match this rule to.
* `group_by` - (List of String) A list of alert labels to group alerts into notifications by.

#### Optional:

* `mute_timings` - (List of String) A list of mute timing names to apply to alerts that match this policy. **Warning** - insert names of mute-timing that already exists, otherwise it can cause problems in your system.
* `continue` - (Boolean) Whether to continue matching subsequent rules if an alert matches the current rule. Otherwise, the rule will be 'consumed' by the first policy to match it.
* `group_wait` - (String) Time to wait to buffer alerts of the same group before sending a notification.
* `group_interval` - (String) Minimum time interval between two notifications for the same group.
* `repeat_interval` - (String) Minimum time interval for re-sending a notification if an alert is still firing.
* `matcher` - (Block List) Describes which labels this rule should match. When multiple matchers are supplied, an alert must match ALL matchers to be accepted by this policy. When no matchers are supplied, the rule will match all alert instances. See below for **nested schema**.
* `policy` - (Block List) Routing rules for specific label sets. See below for **nested schema**.

### Nested schema for `policy.policy.policy.policy.matcher`:

#### Required:

* `label` - (String) The name of the label to match against.
* `match` - (String) The operator to apply when matching values of the given label. Allowed operators are `=` (for equality), `!=` (for negated equality), `=~` (for regex equality), and `!~` (for negated regex equality).
* `value` -  (String) The label value to match against.

## Attribute Reference

* `id` - (String) The ID of this resource. The provider generates the ID, since it's not generated by the API, and it will always be the same ID, since we only use one resource to manage the entire notification policy tree.

### Import Logz.io Grafana Notification Policy as Terraform resource

Since the policies are managed as a tree, and the API itself does not create an ID, we use a constant value for the ID.
You can import existing notification policy as follows:

```
terraform import logzio_grafana_notification_policy.my_np "logzio_policy"
```