package couchbase

import (
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getValidateAllowSaslMechanismDiagMessage function returns pointer to diag message for
// "allow sasl mechanism" parameter
func getValidateAllowSaslMechanismDiagMessage(value string) *diag.Diagnostic {
	return &diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("Allow sasl mechanism type doesn't exist %s\n", value),
		Detail: fmt.Sprintf("Allow sasl mechanism must be:\n%s\n%s\n%s\n%s\n",
			gocb.PlainSaslMechanism,
			gocb.ScramSha1SaslMechanism,
			gocb.ScramSha256SaslMechanism,
			gocb.ScramSha512SaslMechanism,
		),
	}
}

// validateAllowSaslMechanism function verify if allowed sasl mechanism for provider
// Allowed values:
// - PLAIN
// - SCRAM-SHA1
// - SCRAM-SHA256
// - SCRAM-SHA512
func validateAllowSaslMechanism() schema.SchemaValidateDiagFunc {
	return func(i interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		rawSaslMechanism, ok := i.(string)
		if !ok {
			return diag.Errorf("value error: allow sasl mechanism type")
		}

		for _, value := range strings.Split(strings.ReplaceAll(rawSaslMechanism, " ", ""), ",") {
			switch gocb.SaslMechanism(value) {
			case gocb.PlainSaslMechanism,
				gocb.ScramSha1SaslMechanism,
				gocb.ScramSha256SaslMechanism,
				gocb.ScramSha512SaslMechanism:
			default:
				diags = append(diags, *getValidateAllowSaslMechanismDiagMessage(value))
			}
		}
		return diags
	}
}

// validateTLSRootCert function validate TLS root certificate
func validateTLSRootCert() schema.SchemaValidateDiagFunc {
	return func(i interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		tlsRootCAs := *x509.NewCertPool()

		filePath, ok := i.(string)
		if !ok {
			return diag.Errorf("value error: certificate path")
		}

		if filePath != "" {
			file, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
			if err != nil {
				return diag.FromErr(err)
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				return diag.FromErr(err)
			}

			ok = tlsRootCAs.AppendCertsFromPEM(data)
			if !ok {
				return diag.Errorf("cannot append certificate")
			}
		}

		return diags
	}
}
