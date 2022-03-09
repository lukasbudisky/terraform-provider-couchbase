package couchbase

import (
	"context"
	"errors"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: createSecurityUser,
		ReadContext:   readSecurityUser,
		UpdateContext: updateSecurityUser,
		DeleteContext: deleteSecurityUser,
		Description:   "Manage users in couchbase",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			keySecurityUserUsername: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User name",
			},
			keySecurityUserDisplayName: {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "Full user name",
			},
			keySecurityUserPassword: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Sensitive:   true,
				Description: "Password",
			},
			keySecurityUserRole: {
				Type:        schema.TypeSet,
				Elem:        roleStructure(),
				Optional:    true,
				Required:    false,
				ForceNew:    false,
				Description: "User role",
			},
			keySecurityUserGroup: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Required:    false,
				ForceNew:    false,
				Description: "Assigned groups",
			},
		},
	}
}

// userSettings return user settings structure
func userSettings(
	username string,
	displayName string,
	password string,
	rawRoles interface{},
	rawGroups []interface{}) (*gocb.User, error) {

	roles, err := convertRolesToList(rawRoles)
	if err != nil {
		return nil, err
	}

	var groups []string
	for _, group := range rawGroups {
		groups = append(groups, group.(string))
	}

	return &gocb.User{
		Username:    username,
		DisplayName: displayName,
		Groups:      groups,
		Roles:       roles,
		Password:    password,
	}, nil
}

func createSecurityUser(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	us, err := userSettings(
		d.Get(keySecurityUserUsername).(string),
		d.Get(keySecurityUserDisplayName).(string),
		d.Get(keySecurityUserPassword).(string),
		d.Get(keySecurityUserRole),
		d.Get(keySecurityUserGroup).([]interface{}),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = couchbase.UserManager.UpsertUser(*us, nil); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(us.Username)

	return readSecurityUser(c, d, m)
}

func readSecurityUser(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	userID := d.Id()

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	user, err := couchbase.UserManager.GetUser(userID, nil)
	if err != nil && errors.Is(err, gocb.ErrUserNotFound) {
		d.SetId("")
		return diags
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(keySecurityUserUsername, user.Username); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityUserUsername, user.Username, err))
	}

	if err := d.Set(keySecurityUserDisplayName, user.DisplayName); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityUserDisplayName, user.DisplayName, err))
	}

	// Skip set password value because we want detect changes only based on terraform state file

	groupSet := convertRolesToSet(user.Roles)

	if err := d.Set(keySecurityUserRole, groupSet); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityUserRole, user.Roles, err))
	}

	if err := d.Set(keySecurityUserGroup, user.Groups); err != nil {
		diags = append(diags, *diagForValueSet(keySecurityUserGroup, user.Groups, err))
	}

	return diags
}

func updateSecurityUser(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	userID := d.Id()

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	if d.HasChanges(
		keySecurityUserUsername,
		keySecurityUserDisplayName,
		keySecurityUserPassword,
		keySecurityUserRole,
		keySecurityUserGroup,
	) {

		us, err := userSettings(
			userID,
			d.Get(keySecurityUserDisplayName).(string),
			d.Get(keySecurityUserPassword).(string),
			d.Get(keySecurityUserRole),
			d.Get(keySecurityUserGroup).([]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := couchbase.UserManager.UpsertUser(*us, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return readSecurityUser(c, d, m)
}

func deleteSecurityUser(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	userID := d.Id()

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	if err := couchbase.UserManager.DropUser(userID, nil); err != nil {
		diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
