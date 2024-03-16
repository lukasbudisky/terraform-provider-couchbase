package couchbase

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
)

// CollectionSettings custom structure for collection configuration
type CollectionSettings struct {
	Name     string
	Bucket   string
	Scope    string
	Settings *gocb.CreateCollectionSettings
}

// ErrCollectionNotFound custom scope error structure
type ErrCollectionNotFound struct {
	name string
}

// Error function returns custom message when collection is not found
func (e *ErrCollectionNotFound) Error() string {
	return fmt.Sprintf("cannot find collection with name: %s", e.name)
}

// findCollection function will find collection based on name and scope name in couchbase.
// custom error message is returnet when scope is not found
func findCollection(cm *gocb.CollectionManagerV2, name string, scopeName string) (*gocb.CollectionSpec, error) {
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

// getCollectionHistorySettings function creates collection history settings struct based on existing bucket storage type.
// It returns nil if bucket storage type is not magma.
func (cc *Configuration) getCollectionHistorySettings(bucketName string, history bool) (*gocb.CollectionHistorySettings, error) {
	bucket, err := cc.BucketManager.GetBucket(bucketName, nil)

	if err != nil {
		return nil, err
	}

	if bucket.StorageBackend == gocb.StorageBackendMagma {
		return &gocb.CollectionHistorySettings{
			Enabled: history,
		}, nil
	}

	return nil, nil
}
