package couchbase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: createBucket,
		ReadContext:   readBucket,
		UpdateContext: updateBucket,
		DeleteContext: deleteBucket,
		Description:   "Manage buckets in couchbase",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			keyBucketName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket name",
			},
			keyBucketFlushEnabled: {
				Type:        schema.TypeBool,
				ForceNew:    false,
				Default:     true,
				Optional:    true,
				Description: "Bucket flush enable/disable",
			},
			keyBucketQuota: {
				Type:        schema.TypeInt,
				ForceNew:    false,
				Required:    true,
				Description: "Ram quota for bucket",
			},
			keyBucketIndexReplicas: {
				Type:        schema.TypeBool,
				Default:     true,
				ForceNew:    true,
				Optional:    true,
				Description: "Bucket index replicas",
			},
			keyBucketMaxExpiry: {
				Type:        schema.TypeInt,
				Default:     10,
				ForceNew:    false,
				Optional:    true,
				Description: "Max expiry in seconds",
			},
			keyBucketNumReplicas: {
				Type:        schema.TypeInt,
				Default:     1,
				ForceNew:    false,
				Optional:    true,
				Description: "Number of bucket replicas",
			},
			keyBucketBucketType: {
				Type:     schema.TypeString,
				Default:  gocb.CouchbaseBucketType,
				ForceNew: true,
				Optional: true,
				Description: fmt.Sprintf("Bucket type:\n%s\n%s\n%s\n",
					gocb.MemcachedBucketType,
					gocb.EphemeralBucketType,
					gocb.CouchbaseBucketType,
				),
				ValidateDiagFunc: validateBucketType(),
			},
			keyBucketEvictionPolicyType: {
				Type:     schema.TypeString,
				Default:  gocb.EvictionPolicyTypeValueOnly,
				ForceNew: false,
				Optional: true,
				Description: fmt.Sprintf("Eviction policy type:\n%s\n%s\n%s\n%s\n",
					gocb.EvictionPolicyTypeFull,
					gocb.EvictionPolicyTypeValueOnly,
					gocb.EvictionPolicyTypeNotRecentlyUsed,
					gocb.EvictionPolicyTypeNoEviction,
				),
				ValidateDiagFunc: validateEvictionPolicyType(),
			},
			keyBucketCompressionMode: {
				Type:     schema.TypeString,
				Default:  gocb.CompressionModeOff,
				ForceNew: false,
				Optional: true,
				Description: fmt.Sprintf("Compression mode:\n%s\n%s\n%s\n",
					gocb.CompressionModeOff,
					gocb.CompressionModeActive,
					gocb.CompressionModePassive,
				),
				ValidateDiagFunc: validateCompressionMode(),
			},
			keyBucketConflictResolutionType: {
				Type:     schema.TypeString,
				Default:  gocb.ConflictResolutionTypeSequenceNumber,
				ForceNew: true,
				Optional: true,
				Description: fmt.Sprintf("Conflict resolution type:\n%s\n%s\n",
					gocb.ConflictResolutionTypeSequenceNumber,
					gocb.ConflictResolutionTypeTimestamp,
				),
				ValidateDiagFunc: validateConflictResolutionType(),
			},
			keyBucketDurabilityLevel: {
				Type:     schema.TypeInt,
				Default:  gocb.DurabilityLevelNone,
				ForceNew: false,
				Optional: true,
				Description: fmt.Sprintf("Durability level:\n%d\n%d\n%d\n%d\n",
					gocb.DurabilityLevelNone,
					gocb.DurabilityLevelMajority,
					gocb.DurabilityLevelMajorityAndPersistOnMaster,
					gocb.DurabilityLevelPersistToMajority,
				),
				ValidateDiagFunc: validateDurabilityLevel(),
			},
		},
	}
}

func bucketSettings(
	name string,
	flushEnabled bool,
	ramQuota int,
	indexReplicas bool,
	numReplicas int,
	maxExpiry int,
	bucketType string,
	evictionPolicyType string,
	compressionMode string,
	conflictResolutionType string,
	durabilityLevel int) *gocb.CreateBucketSettings {

	return &gocb.CreateBucketSettings{
		BucketSettings: gocb.BucketSettings{
			Name:                   name,
			FlushEnabled:           flushEnabled,
			ReplicaIndexDisabled:   indexReplicas,
			RAMQuotaMB:             uint64(ramQuota),
			NumReplicas:            uint32(numReplicas),
			BucketType:             gocb.BucketType(bucketType),
			EvictionPolicy:         gocb.EvictionPolicyType(evictionPolicyType),
			MaxExpiry:              time.Duration(maxExpiry) * time.Second,
			CompressionMode:        gocb.CompressionMode(compressionMode),
			MinimumDurabilityLevel: gocb.DurabilityLevel(uint8(durabilityLevel)),
		},
		ConflictResolutionType: gocb.ConflictResolutionType(conflictResolutionType),
	}
}

func createBucket(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	bs := bucketSettings(
		d.Get(keyBucketName).(string),
		d.Get(keyBucketFlushEnabled).(bool),
		d.Get(keyBucketQuota).(int),
		d.Get(keyBucketIndexReplicas).(bool),
		d.Get(keyBucketNumReplicas).(int),
		d.Get(keyBucketMaxExpiry).(int),
		d.Get(keyBucketBucketType).(string),
		d.Get(keyBucketEvictionPolicyType).(string),
		d.Get(keyBucketCompressionMode).(string),
		d.Get(keyBucketConflictResolutionType).(string),
		d.Get(keyBucketDurabilityLevel).(int),
	)

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}

	if err := couchbase.BucketManager.CreateBucket(*bs, nil); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bs.Name)

	return readBucket(c, d, m)
}

func readBucket(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	bucketID := d.Id()

	couchbaseConf := m.(*CouchbaseConnection)
	couchbase, diags := couchbaseConf.CouchbaseInitialization()
	if diags != nil {
		return diags
	}

	bucket, err := couchbase.BucketManager.GetBucket(bucketID, nil)
	if err != nil && errors.Is(err, gocb.ErrBucketNotFound) {
		d.SetId("")
		return diags
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(keyBucketName, bucket.Name); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketName, bucket.Name, err))
	}

	if err := d.Set(keyBucketFlushEnabled, bucket.FlushEnabled); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketFlushEnabled, bucket.FlushEnabled, err))
	}

	if err := d.Set(keyBucketQuota, bucket.RAMQuotaMB); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketQuota, bucket.RAMQuotaMB, err))
	}

	if err := d.Set(keyBucketIndexReplicas, bucket.ReplicaIndexDisabled); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketIndexReplicas, bucket.ReplicaIndexDisabled, err))
	}

	if err := d.Set(keyBucketMaxExpiry, int(time.Duration(bucket.MaxExpiry)/time.Second)); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketMaxExpiry, int(time.Duration(bucket.MaxExpiry)/time.Second), err))
	}

	if err := d.Set(keyBucketNumReplicas, bucket.NumReplicas); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketNumReplicas, bucket.NumReplicas, err))
	}

	if err := d.Set(keyBucketBucketType, bucket.BucketType); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketBucketType, bucket.BucketType, err))
	}

	if err := d.Set(keyBucketEvictionPolicyType, bucket.EvictionPolicy); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketEvictionPolicyType, bucket.EvictionPolicy, err))
	}

	if err := d.Set(keyBucketCompressionMode, bucket.CompressionMode); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketCompressionMode, bucket.CompressionMode, err))
	}

	crt, err := couchbaseConf.getBucketConflictResolutionType(bucket.Name)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary: fmt.Sprintf("cannot download couchbase data for %s \n",
				keyBucketConflictResolutionType),
			Detail: fmt.Sprintf("error details: %s\n", err),
		})
	} else {
		if err := d.Set(keyBucketConflictResolutionType, crt); err != nil {
			diags = append(diags, *diagForValueSet(keyBucketConflictResolutionType, crt, err))
		}
	}

	if err := d.Set(keyBucketDurabilityLevel, bucket.MinimumDurabilityLevel); err != nil {
		diags = append(diags, *diagForValueSet(keyBucketDurabilityLevel, bucket.MinimumDurabilityLevel, err))
	}

	return diags
}

func updateBucket(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	bucketID := d.Id()

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}

	if d.HasChanges(
		keyBucketName,
		keyBucketFlushEnabled,
		keyBucketQuota,
		keyBucketIndexReplicas,
		keyBucketNumReplicas,
		keyBucketMaxExpiry,
		keyBucketBucketType,
		keyBucketEvictionPolicyType,
		keyBucketCompressionMode,
		keyBucketDurabilityLevel,
	) {

		bs := bucketSettings(
			bucketID,
			d.Get(keyBucketFlushEnabled).(bool),
			d.Get(keyBucketQuota).(int),
			d.Get(keyBucketIndexReplicas).(bool),
			d.Get(keyBucketNumReplicas).(int),
			d.Get(keyBucketMaxExpiry).(int),
			d.Get(keyBucketBucketType).(string),
			d.Get(keyBucketEvictionPolicyType).(string),
			d.Get(keyBucketCompressionMode).(string),
			d.Get(keyBucketConflictResolutionType).(string),
			d.Get(keyBucketDurabilityLevel).(int),
		)

		if err := couchbase.BucketManager.UpdateBucket(bs.BucketSettings, nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return readBucket(c, d, m)
}

func deleteBucket(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	bucketID := d.Id()

	couchbase, diags := m.(*CouchbaseConnection).CouchbaseInitialization()
	if diags != nil {
		return diags
	}

	if err := couchbase.BucketManager.DropBucket(bucketID, nil); err != nil {
		diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
