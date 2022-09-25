package couchbase

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/couchbase/gocb/v2"
)

// getBucketConflictResolutionType custom fuction for get bucket conflict resolution type because couchbase golang sdk doesn't support to get conflict
// resolution type in gocb v2 version. Currently gocb v2 doesn't support some operations so we must use http/https
// and client-to-node ports for connection. You can read more about ports here:
// https://docs.couchbase.com/server/current/install/install-ports.html
func (cc *Connection) getBucketConflictResolutionType(bucketName string) (*gocb.ConflictResolutionType, error) {
	var (
		conflictResolutionType conflictResolutionType
		scheme                 string
	)

	if cc.ClusterOptions.SecurityConfig.TLSRootCAs == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	client := http.Client{
		Timeout: cc.ClusterOptions.TimeoutsConfig.ManagementTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cc.ClusterOptions.SecurityConfig.TLSSkipVerify,
				RootCAs:            cc.ClusterOptions.SecurityConfig.TLSRootCAs,
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s:%d/pools/default/buckets/%s", scheme, cc.Address, cc.ClientPort, bucketName), nil)
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
