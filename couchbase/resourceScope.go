package couchbase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: createScope,
		ReadContext:   readScope,
		DeleteContext: deleteScope,
		Description:   "Manage scopes in couchbase",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			keyScopeBucketName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket name",
			},
			keyScopeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Scope name",
			},
		},
	}
}

func scopeSettings(
	name string,
	bucket string,
) *ScopeSettings {
	return &ScopeSettings{
		Name:   name,
		Bucket: bucket,
	}
}

func createScope(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	ss := scopeSettings(
		d.Get(keyScopeName).(string),
		d.Get(keyScopeBucketName).(string),
	)

	cm := couchbase.Cluster.Bucket(ss.Bucket).CollectionsV2()

	if err := cm.CreateScope(ss.Name, nil); err != nil {
		return diag.FromErr(err)
	}

	if err := retry.RetryContext(c, time.Duration(scopeTimeoutCreate)*time.Second, func() *retry.RetryError {

		target := &ErrScopeNotFound{}
		_, err := findScope(cm, ss.Name)
		if errors.As(err, &target) {
			return retry.RetryableError(target)
		}

		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("can't create scope: %s error: %s", ss.Name, err))
		}

		d.SetId(ss.Bucket + "/" + ss.Name)
		return nil
	}); err != nil {
		return diag.FromErr(err)
	}

	return readScope(c, d, m)
}

func deleteScope(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	bucketName, scopeName, found := strings.Cut(d.Id(), "/")
	if !found {
		return diag.Errorf("cannot delete scope due to malformed ID: %s", d.Id())
	}

	cm := couchbase.Cluster.Bucket(bucketName).CollectionsV2()

	if err := cm.DropScope(scopeName, nil); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readScope(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	bucketName, scopeName, found := strings.Cut(d.Id(), "/")
	if !found {
		return diag.Errorf("cannot read scope due to malformed ID: %s", d.Id())
	}

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	cm := couchbase.Cluster.Bucket(bucketName).CollectionsV2()

	scope, err := findScope(cm, scopeName)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set(keyScopeName, scope.Name); err != nil {
		diags = append(diags, *diagForValueSet(keyScopeName, scope.Name, err))
	}
	if err := d.Set(keyScopeBucketName, bucketName); err != nil {
		diags = append(diags, *diagForValueSet(keyScopeBucketName, bucketName, err))
	}

	return diags
}
