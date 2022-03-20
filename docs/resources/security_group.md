---
layout: "couchbase"
page_title: "terraform-provider-couchbase resource: couchbase_security_group"
sidebar_current: "docs-couchbase-resource-couchbase_security_group"
description: |-
  Manage groups in couchbase
---

# couchbase_security_group

The `couchbase_security_group` manage groups in couchbase


## Argument reference

The following arguments are supported
### Required

- **name** (String) Group name

### Optional
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>description</b>  (String) Group description</li>
  <li><b>ldap_reference</b> (String) Group ldap reference</li>
  <li><b>role</b> (Block Set) User role. Read more in couchbase documentation - <a href=https://docs.couchbase.com/server/current/rest-api/rbac.html>Role-Based Access Control (RBAC)</a></li>
    <ul>
      <li><b>required nested parameters</b></li>
      <ul>
        <li><b>name</b> (String) Role name</li>
        <li><b>bucket</b> (String) Bucket name</li>
      </ul>
      <li><b>optional nested parameters</b></li>
      <ul>
        <li><b>scope</b> (String) Scope within a bucket</li>
        <li><b>collection</b> (String) Collection within a scope</li>
      </ul>
    </ul>
</ul>

## Attributes reference
The following arguments are exported
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>name</b> (String) Group name</li>
  <li><b>ldap_reference</b> (String) Group ldap reference</li>
  <li><b>role</b> (List - Block Set) User role</li>
</ul>

## Example usage
```
resource "couchbase_security_group" "user_group_1" {
  name        = "user_group_1"
  description = "user group"

  role {
    name   = "query_update"
    bucket = "*"
  }

  role {
    name       = "query_select"
    bucket     = "*"
    scope      = ""
    collection = ""
  }
}
```

## Import

```
# Format:
# terraform import couchbase_security_group.resource_name group_name

# Import command:
terraform import couchbase_security_group.user_group_1 user_group_1
```

