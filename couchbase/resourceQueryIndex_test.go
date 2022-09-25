package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccQueryIndexExtended = `
resource "couchbase_bucket_manager" "bucket" {
	name                     = "testAccPrimaryQueryIndex_extended_bucket_name"
	ram_quota_mb             = 100
	flush_enabled            = false
	max_expire               = 0
	conflict_resolution_type = "seqno"
	compression_mode         = "passive"
	num_replicas             = 0
}

resource "couchbase_query_index" "query_index" {
	name        = "testAccQueryIndex_extended_query_index_name"
	bucket      = couchbase_bucket_manager.bucket.name
	
	fields      = [
		"` + "`" + "action" + "`" + `"
	]
	
	num_replica = 0
	condition   = "(` + "`" + "type" + "`" + " " + `= \"http://example.com\")"
}
`

// TestAccQueryIndex function verify
// - query index extended configuration
func TestAccQueryIndex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccQueryIndexExtended,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "name", "testAccPrimaryQueryIndex_extended_bucket_name"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "ram_quota_mb", "100"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "flush_enabled", "false"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "max_expire", "0"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "conflict_resolution_type", "seqno"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "compression_mode", "passive"),
					resource.TestCheckResourceAttr("couchbase_bucket_manager.bucket", "num_replicas", "0"),
					resource.TestCheckResourceAttr("couchbase_query_index.query_index", "name", "testAccQueryIndex_extended_query_index_name"),
					resource.TestCheckResourceAttr("couchbase_query_index.query_index", "bucket", "testAccPrimaryQueryIndex_extended_bucket_name"),
					resource.TestCheckResourceAttr("couchbase_query_index.query_index", "fields.0", "`action`"),
					resource.TestCheckResourceAttr("couchbase_query_index.query_index", "num_replica", "0"),
					resource.TestCheckResourceAttr("couchbase_query_index.query_index", "condition", "(`type` = \"http://example.com\")"),
				),
			},
		},
	})
}
