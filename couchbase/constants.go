package couchbase

const (
	// Terraform
	terraformVersion = "0.14.9"

	// Provider variables
	providerAddress               = "address"
	providerClientPort            = "client_port"
	providerNodePort              = "node_port"
	providerUsername              = "username"
	providerPassword              = "password"
	providerConnectionTimeout     = "management_timeout"
	providerTLSRootCertSkipVerify = "tls_root_cert_skip_verify"
	providerAllowSaslMechanism    = "allow_sasl_mechanism"
	providerTLSRootCert           = "tls_root_cert"

	// Bucket resource constants, contents
	keyBucketName                   = "name"
	keyBucketFlushEnabled           = "flush_enabled"
	keyBucketQuota                  = "ram_quota_mb"
	keyBucketIndexReplicas          = "replica_index_disable"
	keyBucketMaxExpiry              = "max_expire"
	keyBucketNumReplicas            = "num_replicas"
	keyBucketBucketType             = "bucket_type"
	keyBucketEvictionPolicyType     = "eviction_policy_type"
	keyBucketCompressionMode        = "compression_mode"
	keyBucketConflictResolutionType = "conflict_resolution_type"
	keyBucketDurabilityLevel        = "durability_level"
	keyBucketStorageBackend         = "storage_backend"

	// Security group resource constants, contents
	keySecurityGroupName           = "name"
	keySecurityGroupDescription    = "description"
	keySecurityGroupRole           = "role"
	keySecurityGroupRoleName       = "name"
	keySecurityGroupRoleBucket     = "bucket"
	keySecurityGroupRoleScope      = "scope"
	keySecurityGroupRoleCollection = "collection"
	keySecurityGroupLdapReference  = "ldap_reference"

	// Security user resource constants, contents
	keySecurityUserUsername    = "username"
	keySecurityUserDisplayName = "display_name"
	keySecurityUserPassword    = "password"
	keySecurityUserRole        = "role"
	keySecurityUserGroup       = "groups"

	// Primary query index resource constants, contents
	keyPrimaryQueryIndexName       = "name"
	keyPrimaryQueryIndexBucket     = "bucket"
	keyPrimaryQueryIndexNumReplica = "num_replica"

	// Query index resource constants, contents
	keyQueryIndexName       = "name"
	keyQueryIndexBucket     = "bucket"
	keyQueryIndexNumReplica = "num_replica"
	keyQueryIndexFields     = "fields"
	keyQueryIndexCondition  = "condition"

	// Scope resource constants
	keyScopeName       = "name"
	keyScopeBucketName = "bucket"

	// Collection resource constants
	keyCollectionName       = "name"
	keyCollectionScopeName  = "scope"
	keyCollectionBucketName = "bucket"
	keyCollectionMaxExpiry  = "max_expire"
	keyCollectionHistory    = "history"

	// Others
	queryIndexTimeoutCreate    = 300
	bucketTimeoutCreate        = 300
	scopeTimeoutCreate         = 300
	collectionTimeoutCreate    = 300
	securityUserTimeoutCreate  = 300
	securityGroupTimeoutCreate = 300
)
