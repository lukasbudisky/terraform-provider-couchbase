package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccScopeBasic = `
resource "couchbase_bucket_manager" "bucket" {
    name         = "testAccScope_basic_bucket"
    ram_quota_mb = 100
}

resource "couchbase_bucket_scope" "scope" {
    name   = "testAccScope_basic_scope"
    bucket = couchbase_bucket_manager.bucket.name
}
`

func TestAccScope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScopeBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_bucket_scope.scope", "name", "testAccScope_basic_scope"),
					resource.TestCheckResourceAttr("couchbase_bucket_scope.scope", "bucket", "testAccScope_basic_bucket"),
				),
			},
		},
	})
}
