---
layout: "couchbase"
page_title: "terraform-provider-couchbase resource: couchbase_bucket_manager"
sidebar_current: "docs-couchbase-resource-couchbase_bucket_manager"
description: |-
  Manage buckets in couchbase
---

# couchbase_bucket_manager

The `couchbase_bucket_manager` manage buckets in couchbase


## Argument reference

The following arguments are supported
### Required

- **name** (String) Bucket name
- **ram_quota_mb** (Number) Ram quota for bucket

### Optional
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>bucket_type</b> (String) Bucket type</li>
    <ul>
      <li>memcached</li>
      <li>ephemeral</li>
      <li>membase</li>
    </ul>
  <li><b>compression_mode</b> (String) Compression mode</li>
    <ul>
      <li>off</li>
      <li>active</li>
      <li>passive</li>
    </ul>
  <li><b>conflict_resolution_type</b> (String) Conflict resolution type</li>
    <ul>
      <li>seqno</li>
      <li>lww</li>
    </ul>
  <li><b>durability_level</b> (Number) Durability level</li>
    <ul>
      <li>1</li>
      <li>2</li>
      <li>3</li>
      <li>4</li>
    </ul>
  <li><b>eviction_policy_type</b> (String) Eviction policy type</li>
    <ul>
      <li>fullEviction</li>
      <li>valueOnly</li>
      <li>nruEviction</li>
      <li>noEviction</li>
    </ul>
  <li><b>storage_backend</b> (String) Storage backend type</li>
    <ul>
      <li>couchstore</li>
      <li>magma</li>
    </ul>
  <li><b>flush_enabled</b> (Boolean) Bucket flush enable/disable</li>
  <li><b>max_expire</b> (Int) Max expiry in seconds</li>
  <li><b>num_replicas</b> (Int) Number of bucket replicas</li>
  <li><b>replica_index_disable</b> (Boolean) Bucket index replicas</li>
</ul>

## Attributes reference
The following arguments are exported
<ul>
  <li><b>id</b> (String) The ID of this resource</li>
  <li><b>bucket_name</b> (String) Bucket name</li>
  <li><b>ram_quota</b> (Int) Ram quota for bucket</li>
  <li><b>bucket_type</b> (String) Bucket type</li>
  <li><b>compression_mode</b> (String) Compression mode</li>
  <li><b>conflict_resolution_type</b> (String) Conflict resolution type</li>
  <li><b>durability_level</b> (Int) Durability level</li>
  <li><b>eviction_policy_type</b> (String) Eviction policy type</li>
  <li><b>flush_enabled</b> (Boolean) Bucket flush enable/disable</li>
  <li><b>max_expire</b> (Int) Max expiry in seconds</li>
  <li><b>num_replicas</b> (Int) Number of bucket replicas</li>
  <li><b>replica_index_disable</b> (Boolean) Bucket index replicas</li>
</ul>

## Example usage
```
resource "couchbase_bucket_manager" "bucket_1" {
  name                     = "bucket_1"
  ram_quota_mb             = 512
  flush_enabled            = false
  max_expire               = 0
  conflict_resolution_type = "seqno"
  compression_mode         = "passive"
  num_replicas             = 1
}
```

## Import

```
# Format:
# terraform import couchbase_bucket_manager.resource_name bucket_name

# Import command:
terraform import couchbase_bucket_manager.bucket_1 bucket_1
```

