package couchbase

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateRoleParameter function verify if there is no unexpected value in role parameters:
// - scope
// - collection
func validateRoleParameter() schema.SchemaValidateDiagFunc {
	return func(i interface{}, c cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value, ok := i.(string)
		if !ok {
			return diag.Errorf("value error: scope/collection type")
		}

		if value == "*" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Value scope/collection can't be \"%s\"\n", i),
			})
		}
		return diags
	}
}
