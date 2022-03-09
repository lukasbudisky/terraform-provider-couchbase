package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccUserConfig_basic = `
resource "couchbase_security_user" "testAccUserConfig_basic" {
	username = "testAccUserConfig_basic_name"
	password = "testAccUserConfig_basic_password"
}
`

const testAccUserConfig_extended = `
resource "couchbase_security_group" "group" {
	name        = "testAccUserConfig_extended_group_name"
	description = "testAccUserConfig_extended_group_description"
}

resource "couchbase_security_user" "testAccUserConfig_basic" {
	username = "testAccUserConfig_extended_username"
	password = "testAccUserConfig_extended_password"
	groups = [couchbase_security_group.group.name]
}
`

// TestAccUser function verify
// - user basic configuration
// - user extended configuration
func TestAccUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "id", "testAccUserConfig_basic_name"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "username", "testAccUserConfig_basic_name"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "password", "testAccUserConfig_basic_password"),
					resource.TestCheckNoResourceAttr("couchbase_security_user.testAccUserConfig_basic", "groups"),
				),
			},
			{
				Config: testAccUserConfig_extended,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_security_group.group", "id", "testAccUserConfig_extended_group_name"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "name", "testAccUserConfig_extended_group_name"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "description", "testAccUserConfig_extended_group_description"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "id", "testAccUserConfig_extended_username"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "username", "testAccUserConfig_extended_username"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "password", "testAccUserConfig_extended_password"),
					resource.TestCheckResourceAttr("couchbase_security_user.testAccUserConfig_basic", "groups.0", "testAccUserConfig_extended_group_name"),
				),
			},
		},
	})
}
