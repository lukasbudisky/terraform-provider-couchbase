package couchbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccGroupConfigBasic = `
resource "couchbase_security_group" "group" {
	name        = "testAccGroup_basic_name"
	description = "testAccGroup_basic_description"
}
`

const testAccGroupConfigExtended = `
resource "couchbase_security_group" "group" {
	name           = "testAccGroup_extended_name"
	description    = "testAccGroup_extended_description"
	
	ldap_reference = "OU=testAccGroup_extended_ldap" 
	
	role {
		name    = "query_update"
		bucket  = "*"
	}
}
`

// TestAccGroup function verify
// - group basic configuration
// - group extended configuration
func TestAccGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_security_group.group", "id", "testAccGroup_basic_name"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "name", "testAccGroup_basic_name"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "description", "testAccGroup_basic_description"),
				),
			},
			{
				Config: testAccGroupConfigExtended,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase_security_group.group", "id", "testAccGroup_extended_name"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "name", "testAccGroup_extended_name"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "description", "testAccGroup_extended_description"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "ldap_reference", "OU=testAccGroup_extended_ldap"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "role.0.name", "query_update"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "role.0.bucket", "*"),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "role.0.scope", ""),
					resource.TestCheckResourceAttr("couchbase_security_group.group", "role.0.collection", ""),
				),
			},
		},
	})
}
