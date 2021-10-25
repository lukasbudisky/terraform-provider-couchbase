package couchbase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePrimaryQueryIndex() *schema.Resource {
	return &schema.Resource{
		CreateContext: createPrimaryQueryIndex,
		ReadContext:   readPrimaryQueryIndex,
		DeleteContext: deletePrimaryQueryIndex,
		Description:   "Manage primary query indexes in couchbase",
		Importer: &schema.ResourceImporter{
			StateContext: importQueryIndex,
		},
		Schema: map[string]*schema.Schema{
			keyPrimaryQueryIndexName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Primary query index name",
			},
			keyPrimaryQueryIndexBucket: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Primary query index bucket name",
			},
			keyPrimaryQueryIndexNumReplica: {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Description: "Primary query index number of replica",
			},
		},
	}
}

func createPrimaryQueryIndex(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const deferred = true

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	indexName := d.Get(keyPrimaryQueryIndexName).(string)
	bucketName := d.Get(keyPrimaryQueryIndexBucket).(string)
	numReplica := d.Get(keyPrimaryQueryIndexNumReplica).(int)

	if err := couchbase.createPrimaryQueryIndex(indexName, bucketName, deferred, numReplica); err != nil {
		return diag.FromErr(err)
	}

	if err := resource.RetryContext(c, time.Duration(queryIndexTimeoutCreate)*time.Second, func() *resource.RetryError {

		idx, err := couchbase.readQueryIndexByName(indexName, bucketName)
		if err != nil {
			return resource.RetryableError(err)
		}

		if idx.IsPrimary && idx.Name == indexName {
			if idx.State != getDeferredState(deferred) {
				return resource.RetryableError(fmt.Errorf("primary query index: %s bucket: %s creation in progress: %s", indexName, bucketName, idx.State))
			} else {
				d.SetId(idx.Id)
				return nil
			}
		}

		return resource.NonRetryableError(fmt.Errorf("primary query index doesn't exist index: %s bucket: %s", indexName, bucketName))
	}); err != nil {
		return diag.FromErr(err)
	}

	return readPrimaryQueryIndex(c, d, m)
}

func readPrimaryQueryIndex(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	idx, err := couchbase.readQueryIndexById(d.Id())
	if err != nil && errors.Is(err, gocb.ErrIndexNotFound) {
		d.SetId("")
		return diags
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if !idx.IsPrimary {
		return diag.Errorf("Index is not primary: (index_id=%s)", d.Id())
	}

	if err := d.Set(keyPrimaryQueryIndexName, idx.Name); err != nil {
		diags = append(diags, *diagForValueSet(keyPrimaryQueryIndexName, idx.Name, err))
	}
	if err := d.Set(keyPrimaryQueryIndexBucket, idx.KeyspaceId); err != nil {
		diags = append(diags, *diagForValueSet(keyPrimaryQueryIndexBucket, idx.KeyspaceId, err))
	}

	return diags
}

func deletePrimaryQueryIndex(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	indexName := d.Get(keyPrimaryQueryIndexName).(string)
	bucketName := d.Get(keyPrimaryQueryIndexBucket).(string)

	qis := gocb.DropPrimaryQueryIndexOptions{
		IgnoreIfNotExists: true,
		CustomName:        indexName,
	}

	if err := couchbase.QueryIndexManager.DropPrimaryIndex(bucketName, &qis); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
