---
layout: "couchbase"
page_title: "terraform-provider-couchbase resource: couchbase_bucket_collection"
sidebar_current: "docs-couchbase-resource-couchbase_bucket_collection"
description: |-
  Manage bucket collections in couchbase
---

# couchbase_bucket_collection

The `couchbase_bucket_collection` manage bucket collections in couchbase


## Argument reference

The following arguments are supported
### Required

- **name** (String) Collection name
- **scope** (String) Scope name
- **bucket** (String) Bucket name

### Optional
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>max_expire</b> (Int) Max expiry in seconds</li>
  <li><b>history</b> (Boolean) Collection history enable/disable. Bucket must have "magma" storage mode. Always "False" when storage type is not "magma"</li>
</ul>

## Attributes reference
The following arguments are exported
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>name</b> (String) Collection name</li>
  <li><b>scope</b> (String) Scope name</li>
  <li><b>bucket</b> (String) Bucket name</li>
  <li><b>max_expire</b> (Int) Max expiry in seconds</li>
  <li><b>history</b> (Boolean) Collection history enable/disable. Bucket must have "magma" storage mode. Always "False" when storage type is not "magma"</li>
</ul>

## Example usage
```terraform
resource "couchbase_bucket_collection" "collection_1" {
  name       = "collection_1"
  scope      = "scope_1"
  bucket     = "bucket_1"
  max_expire = 20
  history    = false
}
```

## Import

```bash
# Format:
# terraform import couchbase_bucket_collection.resource_name bucket_name/scope_name/collection_name

# Import command:
terraform import couchbase_bucket_collection.collection_1 bucket_1/scope_1/collection_1
```

