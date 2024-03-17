package couchbase

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
)

// ScopeSettings custom structure for scope configuration
type ScopeSettings struct {
	Name   string
	Bucket string
}

// ErrScopeNotFound custom scope error structure
type ErrScopeNotFound struct {
	name string
}

// Error function returns custom message when scope is not found
func (e *ErrScopeNotFound) Error() string {
	return fmt.Sprintf("cannot find scope with name: %s", e.name)
}

// findScope function will find scope based on name in couchbase.
// custom error message is returnet when scope is not found
func findScope(cm *gocb.CollectionManagerV2, name string) (*gocb.ScopeSpec, error) {
	scopes, err := cm.GetAllScopes(nil)
	if err != nil {
		return nil, err
	}

	for _, scope := range scopes {
		if scope.Name == name {
			return &scope, nil
		}
	}

	return nil, &ErrScopeNotFound{name: name}
}
