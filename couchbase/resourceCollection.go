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

type ErrCollectionNotFound struct {
	name string
}

func (e *ErrCollectionNotFound) Error() string {
	return fmt.Sprintf("cannot find collection with name: %s", e.name)
}

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
		},
	}
}

func createCollection(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	couchbase, diags := m.(*Connection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}
	defer couchbase.ConnectionCLose()

	bucketName := d.Get(keyCollectionBucketName).(string)
	scopeName := d.Get(keyCollectionScopeName).(string)
	collectionName := d.Get(keyCollectionName).(string)

	cm := couchbase.Cluster.Bucket(bucketName).Collections()

	collectionSpec := gocb.CollectionSpec{Name: collectionName, ScopeName: scopeName}
	if err := cm.CreateCollection(collectionSpec, nil); err != nil {
		return diag.FromErr(err)
	}

	if err := retry.RetryContext(c, time.Duration(collectionTimeoutCreate)*time.Second, func() *retry.RetryError {

		target := &ErrCollectionNotFound{}
		_, err := findCollection(cm, collectionName, scopeName)
		if errors.As(err, &target) {
			return retry.RetryableError(target)
		}

		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("can't create collection: %s error: %s", collectionName, err))
		}

		d.SetId(bucketName + "/" + scopeName + "/" + collectionName)
		return nil
	}); err != nil {
		return diag.FromErr(err)
	}

	return readCollection(c, d, m)
}

func readCollection(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	cm := couchbase.Cluster.Bucket(bucketName).Collections()

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

func deleteCollection(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	cm := couchbase.Cluster.Bucket(bucketName).Collections()

	collectionSpec := gocb.CollectionSpec{Name: collectionName, ScopeName: scopeName}
	if err := cm.DropCollection(collectionSpec, nil); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func findCollection(cm *gocb.CollectionManager, name string, scopeName string) (*gocb.CollectionSpec, error) {
	scope, err := findScope(cm, scopeName)
	if err != nil {
		return nil, err
	}

	for _, collection := range scope.Collections {
		if collection.Name == name {
			return &collection, nil
		}
	}

	return nil, &ErrCollectionNotFound{name: name}
}
