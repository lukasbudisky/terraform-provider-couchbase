package couchbase

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateBucketType function validate bucket type
// - membase
// - memcached
// - ephemeral
func validateBucketType() schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value, ok := i.(string)
		if !ok {
			return diag.Errorf("value error: bucket type")
		}

		switch gocb.BucketType(value) {
		case gocb.EphemeralBucketType,
			gocb.MemcachedBucketType,
			gocb.CouchbaseBucketType:
			break
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Bucket type doesn't exist %s\n", i),
				Detail: fmt.Sprintf("Bucket type must be:\n%s\n%s\n%s\n",
					gocb.MemcachedBucketType,
					gocb.EphemeralBucketType,
					gocb.CouchbaseBucketType,
				),
			})
		}
		return diags
	}
}

// validateEvictionPolicyType function verify bucket eviction policy type
// Allowed values:
// - fullEviction
// - valueOnly
// - nruEviction
// - noEviction
func validateEvictionPolicyType() schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value, ok := i.(string)
		if !ok {
			return diag.Errorf("value error: eviction policy type ")
		}

		switch gocb.EvictionPolicyType(value) {
		case gocb.EvictionPolicyTypeFull,
			gocb.EvictionPolicyTypeValueOnly,
			gocb.EvictionPolicyTypeNotRecentlyUsed,
			gocb.EvictionPolicyTypeNoEviction:
			break
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Eviction policy type doesn't exist %s\n", i),
				Detail: fmt.Sprintf("Eviction policy type must be:\n%s\n%s\n%s\n%s\n",
					gocb.EvictionPolicyTypeFull,
					gocb.EvictionPolicyTypeValueOnly,
					gocb.EvictionPolicyTypeNotRecentlyUsed,
					gocb.EvictionPolicyTypeNoEviction,
				),
			})
		}
		return diags
	}
}

// validateCompressionMode function verify bucket compression mode
// Allowed values:
// - off
// - passive
// - active
func validateCompressionMode() schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value, ok := i.(string)
		if !ok {
			return diag.Errorf("value error: compression mode")
		}

		switch gocb.CompressionMode(value) {
		case gocb.CompressionModeOff,
			gocb.CompressionModeActive,
			gocb.CompressionModePassive:
			break
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Compression mode doesn't exist %s\n", i),
				Detail: fmt.Sprintf("Compression mode must be:\n%s\n%s\n%s\n",
					gocb.CompressionModeOff,
					gocb.CompressionModeActive,
					gocb.CompressionModePassive,
				),
			})
		}
		return diags
	}
}

// validateConflictResolutionType function verify bucket conflict resolution type
// Allowed values:
// - lww
// - seqno
func validateConflictResolutionType() schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value, ok := i.(string)
		if !ok {
			return diag.Errorf("value error: conflict resolution type")
		}

		switch gocb.ConflictResolutionType(value) {
		case gocb.ConflictResolutionTypeSequenceNumber,
			gocb.ConflictResolutionTypeTimestamp:
			break
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Conflict resolution type doesn't exist %s\n", i),
				Detail: fmt.Sprintf("Conflict resolution type must be:\n%s\n%s\n",
					gocb.ConflictResolutionTypeSequenceNumber,
					gocb.ConflictResolutionTypeTimestamp,
				),
			})
		}
		return diags
	}
}

// validateDurabilityLevel function verify bucket durability level
// Allowed values:
// - none
// - majority
// - majorityAndPersistActive
// - persistToMajority
func validateDurabilityLevel() schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value, ok := i.(int)
		if !ok {
			return diag.Errorf("value error: durability level")
		}

		switch gocb.DurabilityLevel(uint8(value)) {
		case
			gocb.DurabilityLevelNone,
			gocb.DurabilityLevelMajority,
			gocb.DurabilityLevelMajorityAndPersistOnMaster,
			gocb.DurabilityLevelPersistToMajority:
			break
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Durability level doesn't exist %s\n", i),
				Detail: fmt.Sprintf("Durability level must be:\n%d\n%d\n%d\n%d",
					gocb.DurabilityLevelNone,
					gocb.DurabilityLevelMajority,
					gocb.DurabilityLevelMajorityAndPersistOnMaster,
					gocb.DurabilityLevelPersistToMajority,
				),
			})
		}
		return diags
	}
}
