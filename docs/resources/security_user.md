---
layout: "couchbase"
page_title: "terraform-provider-couchbase resource: couchbase_security_user"
sidebar_current: "docs-couchbase-resource-couchbase_security_user"
description: |-
  Manage users in couchbase
---

# couchbase_security_user

The `couchbase_security_user` manage users in couchbase


## Argument reference

The following arguments are supported
### Required

- **username** (String) Username
- **password** (String, Sensitive) Password

### Optional
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>display_name</b>  (String) Full username</li>
  <li><b>groups</b> (List of String) Assigned groups</li>
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
  <li><b>username</b> (String) Username</li>
  <li><b>display_name</b> (String) Full username</li>
  <li><b>password</b> (String) Password</li>
  <li><b>role</b> (List - Block Set) User role</li>
  <li><b>groups</b> (List of String) Assigned groups</li>
</ul>

## Example usage
```terraform
resource "random_password" "user_password" {
  length  = 10
  special = false
  lower   = true
  upper   = true
}

resource "couchbase_security_user" "user_1" {
  username = "user_1"
  password = random_password.user_password.result

  groups = "user_group_1"

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

> **WARNING**
> If you want import existing user with existing password configure parameter "password" to empty string. There will be no diff during terraform plan.
>
>For Example:
>```terraform
>"couchbase_security_user"">resource "couchbase_security_user" "user_1" {
>  username = "user_1"
>  password = ""
>
>  groups = [couchbase_security_group.user_group_1.id]
>}
>```

```bash
# Format:
# terraform import couchbase_security_user.resource_name user_name

# Import command:
terraform import couchbase_security_user.user_1 user_1
```

