package couchbase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// diagForValueSet function which create custom diagnostic message with severity error
func diagForValueSet(key string, value interface{}, err error) *diag.Diagnostic {
	return &diag.Diagnostic{
		Severity: diag.Error,
		Summary: fmt.Sprintf("cannot set %s with value %s\n",
			key,
			value),
		Detail: fmt.Sprintf("error details: %s\n", err),
	}
}
