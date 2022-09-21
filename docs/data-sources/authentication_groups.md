# Authentication Groups Datasource

Provides a Logz.io authentication groups resource.

* Learn more about authentication groups in the [Logz.io Docs](https://docs.logz.io/api/#tag/Authentication-groups)

**Note**: For this datasource, you don't need to indicate ID or other identifier.

You'll need to create an empty datasource, and it will be populated by all authentication groups in your Logz.io account. 

## Attribute Reference

* `authentication_group` - (Block List) Details for the authentication groups.

#### Nested schema for `authentication_group`:

* `group` - (String) Name of authentication group.
* `user_role` - (String) User role for that group.
