package couchbase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ErrScopeNotFound struct {
	name string
}

func (e *ErrScopeNotFound) Error() string {
	return fmt.Sprintf("cannot find scope with name: %s", e.name)
}

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

func createScope(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	bucketName := d.Get(keyScopeBucketName).(string)
	scopeName := d.Get(keyScopeName).(string)

	collections := couchbase.Cluster.Bucket(bucketName).Collections()

	if err := collections.CreateScope(scopeName, nil); err != nil {
		return diag.FromErr(err)
	}

	if err := retry.RetryContext(c, time.Duration(scopeTimeoutCreate)*time.Second, func() *retry.RetryError {

		target := &ErrScopeNotFound{}
		_, err := findScope(collections, scopeName)
		if errors.As(err, &target) {
			return retry.RetryableError(target)
		}

		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("can't create scope: %s error: %s", scopeName, err))
		}

		d.SetId(bucketName + "/" + scopeName)
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

	collections := couchbase.Cluster.Bucket(bucketName).Collections()

	if err := collections.DropScope(scopeName, nil); err != nil {
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

	collections := couchbase.Cluster.Bucket(bucketName).Collections()

	scope, err := findScope(collections, scopeName)
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

func findScope(cm *gocb.CollectionManager, name string) (*gocb.ScopeSpec, error) {
	scopes, err := cm.GetAllScopes(nil)
	if err != nil {
		return nil, err
	}

	for _, scope := range scopes {
		if scope.Name == name {
			return &scope, nil
		}
	}

	return nil, &ErrScopeNotFound{name: name}
}
