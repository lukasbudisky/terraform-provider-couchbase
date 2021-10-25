package couchbase

import (
	"context"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		TerraformVersion: terraformVersion,
		Schema: map[string]*schema.Schema{
			providerAddress: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_ADDRESS", ""),
			},
			providerPort: {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_PORT", ""),
			},
			providerUsername: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_USERNAME", ""),
			},
			providerPassword: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_PASSWORD", ""),
				Sensitive:   true,
			},
			providerConnectionTimeout: {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Default:     15,
				DefaultFunc: schema.EnvDefaultFunc("CB_MANAGEMENT_TIMEOUT", ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"couchbase_bucket_manager":      resourceBucket(),
			"couchbase_security_group":      resourceSecurityGroup(),
			"couchbase_security_user":       resourceSecurityUser(),
			"couchbase_primary_query_index": resourcePrimaryQueryIndex(),
			"couchbase_query_index":         resourceQueryIndex(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	cc := &CouchbaseConnection{
		Address: d.Get(providerAddress).(string),
		Port:    d.Get(providerPort).(int),
		ClusterOptions: gocb.ClusterOptions{
			Username: d.Get(providerUsername).(string),
			Password: d.Get(providerPassword).(string),
			TimeoutsConfig: gocb.TimeoutsConfig{
				ManagementTimeout: time.Duration(d.Get(providerConnectionTimeout).(int)) * time.Second,
			},
			// TODO
			SecurityConfig: gocb.SecurityConfig{
				TLSSkipVerify: false,
				AllowedSaslMechanisms: []gocb.SaslMechanism{
					gocb.PlainSaslMechanism,
				},
			},
		},
	}

	_, diags := cc.ConnectionValidate()
	if diags != nil {
		return nil, diags
	}

	return cc, diags
}
