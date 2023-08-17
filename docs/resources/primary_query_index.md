---
layout: "couchbase"
page_title: "terraform-provider-couchbase resource: couchbase_primary_query_index"
sidebar_current: "docs-couchbase-resource-couchbase_primary_query_index"
description: |-
  Manage primary query indexex in couchbase
---

# couchbase_primary_query_index

The `couchbase_primary_query_index` manage primary query indexes in couchbase


## Argument reference

The following arguments are supported
### Required

- **name** (String) Primary query index name
- **bucket** (String) Primary query index bucket name

### Optional
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>num_replica</b> (Int) Number of primary query index replicas</li>
</ul>

## Attributes reference
The following arguments are exported
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>name</b> (String) Primary query index name</li>
  <li><b>bucket</b> (String) Primary query index bucket name</li>
  <li><b>num_replica</b> (Int) Number of primary query index replicas</li>
</ul>

## Example usage
```terraform
resource "couchbase_primary_query_index" "primary_index_1" {
  name   = "primary_index_1"
  bucket = "bucket_1"
}
```

## Import

```bash
# Format:
# terraform import couchbase_primary_query_index.resource_name ID

# Import command:
terraform import couchbase_primary_query_index.primary_index_1 ID
```
