package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCollectionBasic = `
resource "couchbase_bucket_manager" "bucket" {
	name         = "testAccCollection_basic_bucket"
	ram_quota_mb = 100
}

resource "couchbase_bucket_scope" "scope" {
	name   = "testAccCollection_basic_scope"
	bucket = couchbase_bucket_manager.bucket.name
}

resource "couchbase_bucket_collection" "collection" {
	name   = "testAccCollection_basic_bucket"
	scope  = couchbase_bucket_scope.scope.name
	bucket = couchbase_bucket_manager.bucket.name
}
`

func TestAccCollection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "name", "testAccCollection_basic_bucket"),
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "bucket", "testAccCollection_basic_bucket"),
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "scope", "testAccCollection_basic_scope"),
				),
			},
		},
	})
}
