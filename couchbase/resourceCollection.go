package couchbase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: createCollection,
		ReadContext:   readCollection,
		DeleteContext: deleteCollection,
		Description:   "Manage collections in couchbase",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			keyCollectionBucketName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket name",
			},
			keyCollectionScopeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Scope name",
			},
			keyCollectionName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Collection name",
			},
			keyCollectionMaxExpiry: {
				Type:        schema.TypeInt,
				Default:     10,
				ForceNew:    true,
				Optional:    true,
				Description: "Max expiry in seconds",
			},
			keyCollectionHistory: {
				Type:        schema.TypeBool,
				Default:     false,
				ForceNew:    true,
				Optional:    true,
				Description: "Collection history enable/disable. Bucket must have \"magma\" storage mode",
			},
		},
	}
}

// collectionSettings return settings structure for collection resource
func collectionSettings(
	name string,
	bucket string,
	scope string,
	maxExpiry int,
	history *gocb.CollectionHistorySettings,
) *CollectionSettings {
	return &CollectionSettings{
		Name:   name,
		Bucket: bucket,
		Scope:  scope,
		Settings: &gocb.CreateCollectionSettings{
			MaxExpiry: time.Duration(maxExpiry) * time.Second,
			History:   history,
		},
	}
}

func createCollection(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	bucketName := d.Get(keyCollectionBucketName).(string)
	historyInput := d.Get(keyCollectionHistory).(bool)

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	history, err := couchbase.getCollectionHistorySettings(bucketName, historyInput)
	if err != nil {
		return diag.FromErr(err)
	}

	cs := collectionSettings(
		d.Get(keyCollectionName).(string),
		bucketName,
		d.Get(keyCollectionScopeName).(string),
		d.Get(keyCollectionMaxExpiry).(int),
		history,
	)

	cm := couchbase.Cluster.Bucket(cs.Bucket).CollectionsV2()

	if err := cm.CreateCollection(cs.Scope, cs.Name, cs.Settings, nil); err != nil {
		return diag.FromErr(err)
	}

	if err := retry.RetryContext(c, time.Duration(collectionTimeoutCreate)*time.Second, func() *retry.RetryError {

		target := &ErrCollectionNotFound{}
		_, err := findCollection(cm, cs.Name, cs.Scope)
		if errors.As(err, &target) {
			return retry.RetryableError(target)
		}

		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("can't create collection: %s error: %s", cs.Name, err))
		}

		d.SetId(cs.Bucket + "/" + cs.Scope + "/" + cs.Name)
		return nil
	}); err != nil {
		return diag.FromErr(err)
	}

	return readCollection(c, d, m)
}

func readCollection(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	names := strings.Split(d.Id(), "/")
	if len(names) != 3 {
		return diag.Errorf("malformed id for collection: %s", d.Id())
	}

	bucketName := names[0]
	scopeName := names[1]
	collectionName := names[2]

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	cm := couchbase.Cluster.Bucket(bucketName).CollectionsV2()

	collection, err := findCollection(cm, collectionName, scopeName)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set(keyCollectionName, collection.Name); err != nil {
		diags = append(diags, *diagForValueSet(keyCollectionName, collection.Name, err))
	}
	if err := d.Set(keyCollectionBucketName, bucketName); err != nil {
		diags = append(diags, *diagForValueSet(keyCollectionBucketName, bucketName, err))
	}
	if err := d.Set(keyCollectionScopeName, scopeName); err != nil {
		diags = append(diags, *diagForValueSet(keyCollectionScopeName, scopeName, err))
	}

	return diags
}

func deleteCollection(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	names := strings.Split(d.Id(), "/")
	bucketName := names[0]
	scopeName := names[1]
	collectionName := names[2]

	cm := couchbase.Cluster.Bucket(bucketName).CollectionsV2()

	if err := cm.DropCollection(scopeName, collectionName, nil); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
