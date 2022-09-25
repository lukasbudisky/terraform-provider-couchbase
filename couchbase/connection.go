package couchbase

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Connection struct contain information about connection parameters.
type Connection struct {
	Scheme         string
	Address        string
	NodePort       int
	ClientPort     int
	ClusterOptions gocb.ClusterOptions
}

// conflictResolutionType custom struct for bucket conflict resolution type because couchbase golang sdk doesn't support to get conflict
// resolution type in gocb v2 version
type conflictResolutionType struct {
	ConflictResolutionType gocb.ConflictResolutionType `json:"conflictResolutionType"`
}

// Configuration struct contains information about cluster and bucket manager.
type Configuration struct {
	Cluster           *gocb.Cluster
	BucketManager     *gocb.BucketManager
	UserManager       *gocb.UserManager
	QueryIndexManager *gocb.QueryIndexManager
}

// CouchbaseInitialization function creates connection to couchbase.
func (cc *Connection) CouchbaseInitialization() (*Configuration, diag.Diagnostics) {

	cluster, diags := cc.ConnectionValidate()

	return &Configuration{
		Cluster:           cluster,
		BucketManager:     cluster.Buckets(),
		UserManager:       cluster.Users(),
		QueryIndexManager: cluster.QueryIndexes(),
	}, diags
}

// ConnectionValidate function validates connection to couchbase
func (cc *Connection) ConnectionValidate() (*gocb.Cluster, diag.Diagnostics) {
	var diags diag.Diagnostics

	cbAddress := fmt.Sprintf("%s://%s:%d", cc.Scheme, cc.Address, cc.NodePort)

	cluster, err := gocb.Connect(cbAddress, cc.ClusterOptions)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("cannot connect to couchbase %s\n", cbAddress),
			Detail:   fmt.Sprintf("error details: %s\n", err),
		})
		return nil, diags
	}
	return cluster, diags
}

// ConnectionCLose close couchbase connection
func (cc *Configuration) ConnectionCLose() {
	cc.Cluster.Close(nil)
}
