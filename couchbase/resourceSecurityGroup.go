package couchbase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: createSecurityGroup,
		ReadContext:   readSecurityGroup,
		UpdateContext: updateSecurityGroup,
		DeleteContext: deleteSecurityGroup,
		Description:   "Manage groups in couchbase",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			keySecurityGroupName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Group name",
			},
			keySecurityGroupDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Group description",
			},
			keySecurityGroupRole: {
				Type:        schema.TypeSet,
				Elem:        roleStructure(),
				Optional:    true,
				Required:    false,
				ForceNew:    false,
				Description: "Group role",
			},
			keySecurityGroupLdapReference: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Group ldap reference",
			},
		},
	}
}

// groupSettings return group settings structure
func groupSettings(
	name string,
	description string,
	rawRoles interface{},
	ldapGroupReference string) (*gocb.Group, error) {

	roles, err := convertRolesToList(rawRoles)
	if err != nil {
		return nil, err
	}

	return &gocb.Group{
		Name:               name,
		Description:        description,
		Roles:              roles,
		LDAPGroupReference: ldapGroupReference,
	}, nil
}

func createSecurityGroup(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	gs, err := groupSettings(
		d.Get(keySecurityGroupName).(string),
		d.Get(keySecurityGroupDescription).(string),
		d.Get(keySecurityGroupRole),
		d.Get(keySecurityGroupLdapReference).(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = couchbase.UserManager.UpsertGroup(*gs, nil); err != nil {
		return diag.FromErr(err)
	}

	if err := retry.RetryContext(c, time.Duration(securityGroupTimeoutCreate)*time.Second, func() *retry.RetryError {

		_, err := couchbase.UserManager.GetGroup(gs.Name, nil)
		if err != nil && errors.Is(err, gocb.ErrGroupNotFound) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("can't create security group: %s error: %s", gs.Name, err))
		}

		d.SetId(gs.Name)
		return nil
	}); err != nil {
		return diag.FromErr(err)
	}

	return readSecurityGroup(c, d, m)
}

func readSecurityGroup(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	groupID := d.Id()

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	group, err := couchbase.UserManager.GetGroup(groupID, nil)
	if err != nil && errors.Is(err, gocb.ErrGroupNotFound) {
		d.SetId("")
		return diags
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(keySecurityGroupName, group.Name); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityGroupName, group.Name, err))
	}

	if err := d.Set(keySecurityGroupDescription, group.Description); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityGroupDescription, group.Description, err))
	}

	groupSet := convertRolesToSet(group.Roles)

	if err := d.Set(keySecurityGroupRole, groupSet); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityGroupRole, group.Roles, err))
	}

	if err := d.Set(keySecurityGroupLdapReference, group.LDAPGroupReference); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityGroupLdapReference, group.LDAPGroupReference, err))
	}

	return diags
}

func updateSecurityGroup(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	groupID := d.Id()

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	if d.HasChanges(
		keySecurityGroupName,
		keySecurityGroupDescription,
		keySecurityGroupRole,
		keySecurityGroupLdapReference,
	) {

		gs, err := groupSettings(
			groupID,
			d.Get(keySecurityGroupDescription).(string),
			d.Get(keySecurityGroupRole),
			d.Get(keySecurityGroupLdapReference).(string),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := couchbase.UserManager.UpsertGroup(*gs, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return readSecurityGroup(c, d, m)
}

func deleteSecurityGroup(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	groupID := d.Id()

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	if err := couchbase.UserManager.DropGroup(groupID, nil); err != nil {
		diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
