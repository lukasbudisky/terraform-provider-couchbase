package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPrimaryQueryIndexBasic = `
resource "couchbase_bucket_manager" "bucket" {
    name                     = "testAccPrimaryQueryIndex_basic_bucket_name"
    ram_quota_mb             = 100
    flush_enabled            = false
    max_expire               = 0
    conflict_resolution_type = "seqno"
    compression_mode         = "passive"
    num_replicas             = 0
}

resource "couchbase_primary_query_index" "primary_index" {
    name   = "testAccPrimaryQueryIndex_basic_primary_index_name"
    bucket = couchbase_bucket_manager.bucket.name
}
`

// TestAccPrimaryQueryIndex function verify
// - basic primary index query configuration
func TestAccPrimaryQueryIndex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrimaryQueryIndexBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "name", "testAccPrimaryQueryIndex_basic_bucket_name"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "ram_quota_mb", "100"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "flush_enabled", "false"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "max_expire", "0"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "conflict_resolution_type", "seqno"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "compression_mode", "passive"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "num_replicas", "0"),
					resource.TestCheckResourceAttr("couchbase_primary_query_index.primary_index", "name", "testAccPrimaryQueryIndex_basic_primary_index_name"),
					resource.TestCheckResourceAttr("couchbase_primary_query_index.primary_index", "bucket", "testAccPrimaryQueryIndex_basic_bucket_name"),
				),
			},
		},
	})
}
