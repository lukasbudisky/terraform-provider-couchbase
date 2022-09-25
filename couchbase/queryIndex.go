package couchbase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// queryIndex custom query index structure.
// gocb v2 doesn't have ID in structure
type queryIndex struct {
	Condition   string   `json:"condition"`
	DatastoreID string   `json:"datastore_id"`
	ID          string   `json:"id"`
	IndexKey    []string `json:"index_key"`
	IsPrimary   bool     `json:"is_primary"`
	KeyspaceID  string   `json:"keyspace_id"`
	Name        string   `json:"name"`
	NamespaceID string   `json:"namespace_id"`
	State       string   `json:"state"`
	Using       string   `json:"using"`
}

// convertFiedsToList function convert list of fields interfaces to list of strings
func convertFieldsToList(rawFields []interface{}) ([]string, error) {
	var fields []string

	for _, field := range rawFields {
		sub, ok := field.(string)
		if !ok {
			return nil, fmt.Errorf("cannot convert query index field")
		}
		fields = append(fields, sub)
	}
	return fields, nil
}

// getDeferredState function return string based on deferred bool value
func getDeferredState(state bool) string {
	if state {
		return "deferred"
	}
	return "online"
}

// readQueryIndexByID function read query indexes based on ID
func (cc *Configuration) readQueryIndexByID(id string) (*queryIndex, error) {
	q := "SELECT `indexes`.* FROM system:indexes WHERE id=? AND `using`=\"gsi\""
	rows, err := cc.Cluster.Query(q, &gocb.QueryOptions{
		PositionalParameters: []interface{}{id},
		Readonly:             true,
	})
	if err != nil {
		return nil, err
	}

	var index *queryIndex

	for rows.Next() {
		err := rows.Row(&index)
		if err != nil {
			return nil, err
		}
		break
	}
	defer rows.Close()
	if index == nil {
		return nil, fmt.Errorf("index not found id: %s; %w", id, gocb.ErrIndexNotFound)
	}
	return index, nil
}

// readQueryIndexByName function read query indexes based on index name and bucket name
func (cc *Configuration) readQueryIndexByName(indexName, bucketName string) (*queryIndex, error) {
	q := "SELECT `indexes`.* FROM system:indexes WHERE keyspace_id=? AND name=? AND `using`=\"gsi\""
	rows, err := cc.Cluster.Query(q, &gocb.QueryOptions{
		PositionalParameters: []interface{}{bucketName, indexName},
		Readonly:             true,
	})
	if err != nil {
		return nil, err
	}

	var index *queryIndex

	for rows.Next() {
		err := rows.Row(&index)
		if err != nil {
			return nil, err
		}
		break
	}
	defer rows.Close()
	if index == nil {
		return nil, fmt.Errorf("index not found index: %s bucket:%s; %w", indexName, bucketName, gocb.ErrIndexNotFound)
	}
	return index, nil
}

// createPrimaryQueryIndex custom functon which support primary query index creation with deferred state, number of replicas
func (cc *Configuration) createPrimaryQueryIndex(indexName, bucketName string, deferred bool, numReplica int) error {
	q := fmt.Sprintf("CREATE PRIMARY INDEX `%s` ON `%s` WITH {\"defer_build\":%t, \"num_replica\":%d}", indexName, bucketName, deferred, numReplica)
	rows, err := cc.Cluster.Query(q, nil)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// createQueryIndex custom functon which support query index creation with fields parameters and conditions, deferred state, number of replicas
func (cc *Configuration) createQueryIndex(indexName, bucketName string, fields []string, condition string, deferred bool, numReplica int) error {
	if len(fields) <= 0 {
		return fmt.Errorf("you must specify at least one field to index")
	}

	if condition != "" {
		condition = fmt.Sprintf("WHERE %s", condition)
	} else {
		condition = ""
	}

	q := fmt.Sprintf("CREATE INDEX `%s` ON `%s`(%s) %s WITH {\"defer_build\":%t, \"num_replica\":%d}", indexName, bucketName, strings.Join(fields, ","), condition, deferred, numReplica)
	rows, err := cc.Cluster.Query(q, nil)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// parseID function which parse id and number of index replicas during import
func parseID(id string) (string, int, error) {
	results := strings.Split(id, ",")
	if len(results) != 2 {
		return "", 0, fmt.Errorf("cannot parse id during import id: %s", id)
	}

	sub, err := strconv.Atoi(results[1])
	if err != nil {
		return "", 0, fmt.Errorf("cannot convert part of id to int id: %s", id)
	}

	return results[0], sub, nil
}

// importQueryIndex custom terraform resource import function
func importQueryIndex(c context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	id, replica, err := parseID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set(keyQueryIndexNumReplica, replica)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
