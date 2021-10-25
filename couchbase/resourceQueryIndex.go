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

func resourceQueryIndex() *schema.Resource {
	return &schema.Resource{
		CreateContext: createQueryIndex,
		ReadContext:   readQueryIndex,
		DeleteContext: deleteQueryIndex,
		Description:   "Manage query indexes in couchbase",
		Importer: &schema.ResourceImporter{
			StateContext: importQueryIndex,
		},
		Schema: map[string]*schema.Schema{
			keyQueryIndexName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Query index name",
			},
			keyQueryIndexBucket: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Query index bucket name",
			},
			keyQueryIndexFields: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Optional:    false,
				ForceNew:    true,
				Description: "Query index fields",
			},
			keyQueryIndexCondition: {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "Query index where statement",
			},
			keyQueryIndexNumReplica: {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Description: "Query index number of replica",
			},
		},
	}
}

func createQueryIndex(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const deferred = true

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}

	indexName := d.Get(keyQueryIndexName).(string)
	bucketName := d.Get(keyQueryIndexBucket).(string)
	numReplica := d.Get(keyQueryIndexNumReplica).(int)
	condition := d.Get(keyQueryIndexCondition).(string)
	fields, err := convertFieldsToList(d.Get(keyQueryIndexFields).([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := couchbase.createQueryIndex(indexName, bucketName, fields, condition, deferred, numReplica); err != nil {
		return diag.FromErr(err)
	}

	if err := resource.RetryContext(c, time.Duration(queryIndexTimeoutCreate)*time.Second, func() *resource.RetryError {

		idx, err := couchbase.readQueryIndexByName(indexName, bucketName)
		if err != nil {
			return resource.RetryableError(err)
		}

		if !idx.IsPrimary && idx.Name == indexName {
			if idx.State != getDeferredState(deferred) {
				return resource.RetryableError(fmt.Errorf("primary query index: %s bucket: %s creation in progress: %s", indexName, bucketName, idx.State))
			} else {
				d.SetId(idx.Id)
				return nil
			}
		}

		return resource.NonRetryableError(fmt.Errorf("query index doesn't exist index: %s bucket: %s", indexName, bucketName))
	}); err != nil {
		return diag.FromErr(err)
	}

	return readQueryIndex(c, d, m)
}

func readQueryIndex(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}

	idx, err := couchbase.readQueryIndexById(d.Id())
	if err != nil && errors.Is(err, gocb.ErrIndexNotFound) {
		d.SetId("")
		return diags
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if idx.IsPrimary {
		return diag.Errorf("Index is primary: (index_id=%s)", d.Id())
	}

	if err := d.Set(keyPrimaryQueryIndexName, idx.Name); err != nil {
		diags = append(diags, *diagForValueSet(keyPrimaryQueryIndexName, idx.Name, err))
	}
	if err := d.Set(keyPrimaryQueryIndexBucket, idx.KeyspaceId); err != nil {
		diags = append(diags, *diagForValueSet(keyPrimaryQueryIndexBucket, idx.KeyspaceId, err))
	}
	if err := d.Set(keyQueryIndexCondition, idx.Condition); err != nil {
		diags = append(diags, *diagForValueSet(keyQueryIndexCondition, idx.Condition, err))
	}
	if err := d.Set(keyQueryIndexFields, idx.IndexKey); err != nil {
		diags = append(diags, *diagForValueSet(keyQueryIndexFields, idx.IndexKey, err))
	}

	return diags
}

func deleteQueryIndex(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}

	indexName := d.Get(keyQueryIndexName).(string)
	bucketName := d.Get(keyQueryIndexBucket).(string)

	qis := gocb.DropQueryIndexOptions{
		IgnoreIfNotExists: true,
	}

	if err := couchbase.QueryIndexManager.DropIndex(bucketName, indexName, &qis); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
