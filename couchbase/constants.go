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

	// Bucket resource contants
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

	// Security group resource contants
	keySecurityGroupName           = "name"
	keySecurityGroupDescription    = "description"
	keySecurityGroupRole           = "role"
	keySecurityGroupRoleName       = "name"
	keySecurityGroupRoleBucket     = "bucket"
	keySecurityGroupRoleScope      = "scope"
	keySecurityGroupRoleCollection = "collection"
	keySecurityGroupLdapReference  = "ldap_reference"

	// Security user resource contants
	keySecurityUserUsername    = "username"
	keySecurityUserDisplayName = "display_name"
	keySecurityUserPassword    = "password"
	keySecurityUserRole        = "role"
	keySecurityUserGroup       = "groups"

	// Primary query index resource contants
	keyPrimaryQueryIndexName       = "name"
	keyPrimaryQueryIndexBucket     = "bucket"
	keyPrimaryQueryIndexNumReplica = "num_replica"

	// Query index resource contants
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

	// Others
	queryIndexTimeoutCreate    = 300
	bucketTimeoutCreate        = 300
	scopeTimeoutCreate         = 300
	collectionTimeoutCreate    = 300
	securityUserTimeoutCreate  = 300
	securityGroupTimeoutCreate = 300
)
