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

const testAccCollectionBucketMagmaStorage = `
resource "couchbase_bucket_manager" "bucket" {
    name                     = "testAccCollection_bucket_magma_storage"
    ram_quota_mb             = 1024
    bucket_type              = "membase"
    compression_mode         = "passive"
    conflict_resolution_type = "seqno"
    durability_level		 = 1
    eviction_policy_type     = "valueOnly"
    flush_enabled            = false
    max_expire               = 0
    num_replicas             = 0
    replica_index_disable    = true
    storage_backend          = "magma"
}

resource "couchbase_bucket_scope" "scope" {
    name   = "testAccCollection_bucket_magma_scope"
    bucket = couchbase_bucket_manager.bucket.name
}

resource "couchbase_bucket_collection" "collection" {
    name       = "testAccCollection_bucket_magma_collection"
    scope      = couchbase_bucket_scope.scope.name
    bucket     = couchbase_bucket_manager.bucket.name
    max_expire = 20
    history    = true
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

func TestAccCollectionBucketMagmaStorage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionBucketMagmaStorage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "name", "testAccCollection_bucket_magma_collection"),
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "bucket", "testAccCollection_bucket_magma_storage"),
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "scope", "testAccCollection_bucket_magma_scope"),
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "max_expire", "20"),
					resource.TestCheckResourceAttr("couchbase_bucket_collection.collection", "history", "true"),
				),
			},
		},
	})
}
