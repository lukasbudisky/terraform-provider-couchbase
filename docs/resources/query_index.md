---
layout: "couchbase"
page_title: "terraform-provider-couchbase resource: couchbase_query_index"
sidebar_current: "docs-couchbase-resource-couchbase_query_index"
description: |-
  Manage query indexex in couchbase
---

# couchbase_query_index

The `couchbase_query_index` manage query indexes in couchbase


## Argument reference

The following arguments are supported
### Required

- **name** (String) Query index name
- **bucket** (String) Query index bucket name
- **fields** (List of String) Query index fields - This parameter should include also backticks

### Optional
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>num_replica</b> (Int) Number of query index replicas</li>
  <li><b>condition</b> (String) Query index where statement - This parameter should include also backticks<li>
</ul>

## Attributes reference
The following arguments are exported
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>name</b> (String) Query index name</li>
  <li><b>bucket</b> (String) Query index bucket name</li>
  <li><b>condition</b> (String) Query index where statement</li>
  <li><b>fields</b> (List of String) Query index fields</li>
  <li><b>num_replica</b> (Int) Number of query index replicas</li>
</ul>

## Example usage
```
resource "couchbase_query_index" "index_1" {
  name        = "index_1"
  bucket      = "bucket_1"
  fields      = ["`action`"]
  num_replica = 0
  condition   = "(`type` = \"http://example.com\")"
}
```

## Import

```
# Format:
# terraform import couchbase_query_index.resource_name ID

# Import command:
terraform import couchbase_query_index.index_1 ID
```
