package couchbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/couchbase/gocb/v2"
)

// getBucketConflictResolutionType custom fuction for get bucket conflict resolution type because couchbase golang sdk doesn't support to get conflict
// resolution type in gocb v2 version
func (cc *CouchbaseConnection) getBucketConflictResolutionType(bucketName string) (*gocb.ConflictResolutionType, error) {
	var conflictResolutionType conflictResolutionType
	var schema string
	client := http.Client{Timeout: cc.ClusterOptions.TimeoutsConfig.ManagementTimeout}

	// TODO
	if cc.ClusterOptions.SecurityConfig.TLSRootCAs == nil {
		schema = "http"
	} else {
		schema = "https"
	}

	// TODO https
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s:%d/pools/default/buckets/%s", schema, cc.Address, cc.Port, bucketName), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cc.ClusterOptions.Username, cc.ClusterOptions.Password)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resData, &conflictResolutionType); err != nil {
		return nil, err
	}

	return &conflictResolutionType.ConflictResolutionType, nil
}
