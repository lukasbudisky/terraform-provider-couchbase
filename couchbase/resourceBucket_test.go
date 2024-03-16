package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccBucketBasic = `
resource "couchbase_bucket_manager" "bucket" {
    name                     = "testAccBucket_basic_bucket_name"
    ram_quota_mb             = 100
}
`

const testAccBucketExtended = `
resource "couchbase_bucket_manager" "bucket" {
    name                     = "testAccBucket_extended_bucket_name"
    ram_quota_mb             = 1024
    bucket_type              = "membase"
    compression_mode         = "passive"
    conflict_resolution_type = "seqno"
    durability_level         = 1
    eviction_policy_type     = "valueOnly"
    flush_enabled            = false
    max_expire               = 0
    num_replicas             = 0
    replica_index_disable    = true
    storage_backend          = "magma"
}
`

// TestAccBucket function verify
// - basic bucket configuration
// - extended bucket configuration
func TestAccBucket(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "name", "testAccBucket_basic_bucket_name"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "ram_quota_mb", "100"),
				),
			},
			{
				Config: testAccBucketExtended,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "name", "testAccBucket_extended_bucket_name"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "ram_quota_mb", "1024"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "bucket_type", "membase"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "compression_mode", "passive"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "conflict_resolution_type", "seqno"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "durability_level", "1"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "eviction_policy_type", "valueOnly"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "flush_enabled", "false"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "max_expire", "0"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "num_replicas", "0"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "replica_index_disable", "true"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "storage_backend", "magma"),
				),
			},
		},
	})
}
