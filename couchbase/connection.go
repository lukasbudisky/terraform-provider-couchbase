package couchbase

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// CouchbaseConnection struct container information about connection parameters.
type CouchbaseConnection struct {
	Address        string
	Port           int
	ClusterOptions gocb.ClusterOptions
}

// conflictResolutionType custom struct for bucket conflict resolution type because couchbase golang sdk doesn't support to get conflict
// resolution type in gocb v2 version
type conflictResolutionType struct {
	ConflictResolutionType gocb.ConflictResolutionType `json:"conflictResolutionType"`
}

// CouchbaseConfiguration struct contains informatio about cluster and bucket manager.
type CouchbaseConfiguration struct {
	Cluster           *gocb.Cluster
	BucketManager     *gocb.BucketManager
	UserManager       *gocb.UserManager
	QueryIndexManager *gocb.QueryIndexManager
}

// CouchbaseInitialization function creates connection to couchbase.
func (cc *CouchbaseConnection) CouchbaseInitialization() (*CouchbaseConfiguration, diag.Diagnostics) {

	cluster, diags := cc.ConnectionValidate()

	return &CouchbaseConfiguration{
		Cluster:           cluster,
		BucketManager:     cluster.Buckets(),
		UserManager:       cluster.Users(),
		QueryIndexManager: cluster.QueryIndexes(),
	}, diags
}

// ConnectionValidate function validates connection to couchbase
func (cc *CouchbaseConnection) ConnectionValidate() (*gocb.Cluster, diag.Diagnostics) {
	var diags diag.Diagnostics

	cluster, err := gocb.Connect(fmt.Sprintf("%s:%d", cc.Address, cc.Port), cc.ClusterOptions)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("cannot connect to couchbase %s:%d\n", cc.Address, cc.Port),
			Detail:   fmt.Sprintf("error details: %s\n", err),
		})
		return nil, diags
	}
	return cluster, diags
}
