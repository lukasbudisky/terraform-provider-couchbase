package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccUserConfig_basic = `
resource "couchbase_security_user" "testAccUserConfig_basic" {
	username = "testAccUser_basic"
	password = "password"
  }
`

const testAccUserConfig_groups = `
resource "couchbase_security_group" "group" {
	name        = "testAccUser_group"
	description = "user group"
}
resource "couchbase_security_user" "testAccUserConfig_basic" {
	username = "testAccUser_groups"
	password = "password"
	groups = [couchbase_security_group.group.name]
  }
`

func TestAccUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "username", "testAccUser_basic"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "password", "password"),
					resource.TestCheckNoResourceAttr("couchbase_security_user.testAccUserConfig_basic", "groups"),
				),
			},
			{
				Config: testAccUserConfig_groups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_security_group.group", "name", "testAccUser_group"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "description", "user group"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "username", "testAccUser_groups"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "password", "password"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "groups", "testAccUser_group"),
				),
			},
		},
	})
}
