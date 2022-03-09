package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccBucket_basic = `
resource "couchbase_bucket_manager" "bucket" {
	name                     = "testAccPrimaryQueryIndex_basic_bucket_name"
	ram_quota_mb             = 100
}
`

const testAccBucket_extended = `
resource "couchbase_bucket_manager" "bucket" {
	name                     = "testAccPrimaryQueryIndex_extended_bucket_name"
	ram_quota_mb             = 100
	bucket_type			     = "membase"
	compression_mode         = "passive"
	conflict_resolution_type = "seqno"
	durability_level		 = 0
	eviction_policy_type     = "valueOnly"
	flush_enabled            = false
	max_expire               = 0
	num_replicas             = 0
	replica_index_disable    = true
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
				Config: testAccBucket_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "name", "testAccPrimaryQueryIndex_basic_bucket_name"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "ram_quota_mb", "100"),
				),
			},
			{
				Config: testAccBucket_extended,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "name", "testAccPrimaryQueryIndex_extended_bucket_name"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "ram_quota_mb", "100"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "bucket_type", "membase"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "compression_mode", "passive"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "conflict_resolution_type", "seqno"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "durability_level", "0"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "eviction_policy_type", "valueOnly"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "flush_enabled", "false"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "max_expire", "0"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "num_replicas", "0"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "replica_index_disable", "true"),
				),
			},
		},
	})
}
