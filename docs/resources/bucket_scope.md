---
layout: "couchbase"
page_title: "terraform-provider-couchbase resource: couchbase_bucket_scope"
sidebar_current: "docs-couchbase-resource-couchbase_bucket_scope"
description: |-
  Manage bucket scopes in couchbase
---

# couchbase_bucket_scope

The `couchbase_bucket_scope` manage bucket scopes in couchbase


## Argument reference

The following arguments are supported
### Required

- **name** (String) Scope name
- **bucket** (String) Bucket name

### Optional
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
</ul>

## Attributes reference
The following arguments are exported
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>name</b> (String) Scope name</li>
  <li><b>bucket</b> (String) Bucket name</li>
</ul>

## Example usage
```terraform
resource "couchbase_bucket_scope" "scope_1" {
  name   = "scope_1"
  bucket = "bucket_1"
}
```

## Import

```bash
# Format:
# terraform import couchbase_bucket_scope.resource_name bucket_name/scope_name

# Import command:
terraform import couchbase_bucket_scope.scope_1 bucket_1/scope_1
```

