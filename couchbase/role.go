package couchbase

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RoleStructure function provide terraform role resource structure
func roleStructure() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			keySecurityGroupRoleName: {
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
				ForceNew: false,
			},
			keySecurityGroupRoleBucket: {
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
				ForceNew: false,
			},
			keySecurityGroupRoleScope: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ForceNew:         false,
				ValidateDiagFunc: validateRoleParameter(),
			},
			keySecurityGroupRoleCollection: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ForceNew:         false,
				ValidateDiagFunc: validateRoleParameter(),
			},
		},
	}
}

// convertRolesToList function convert raw roles to couchbase list of roles
func convertRolesToList(rawRoles interface{}) ([]gocb.Role, error) {
	roles := []gocb.Role{}

	for _, role := range rawRoles.(*schema.Set).List() {
		sub, ok := role.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot convert input roles to couchbase roles")
		}
		roles = append(roles,
			gocb.Role{
				Name:       sub[keySecurityGroupRoleName].(string),
				Bucket:     sub[keySecurityGroupRoleBucket].(string),
				Scope:      sub[keySecurityGroupRoleScope].(string),
				Collection: sub[keySecurityGroupRoleCollection].(string),
			},
		)
	}
	return roles, nil
}

// convertRolesToSet function convert couchbase list of roles to terraform schema set
func convertRolesToSet(createdRoles []gocb.Role) *schema.Set {
	roles := []interface{}{}
	for _, role := range createdRoles {
		rawRole := make(map[string]interface{})
		rawRole[keySecurityGroupRoleName] = role.Name
		rawRole[keySecurityGroupRoleBucket] = role.Bucket
		rawRole[keySecurityGroupRoleCollection] = role.Collection
		rawRole[keySecurityGroupRoleScope] = role.Scope
		roles = append(roles, rawRole)
	}

	s := schema.NewSet(schema.HashResource(roleStructure()), roles)

	return s
}
