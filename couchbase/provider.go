package couchbase

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		TerraformVersion: terraformVersion,
		Schema: map[string]*schema.Schema{
			providerAddress: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_ADDRESS", ""),
				Description: "Couchbase address (without scheme)",
			},
			providerClientPort: {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_CLIENT_PORT", ""),
				Description: "Couchbase port: client-to-node (use for api calls usually http 8091 or https 18091)",
			},
			providerNodePort: {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_NODE_PORT", ""),
				Description: "Couchbase port: node-to-node (usually for scheme couchbase 11210 or couchbases 112107)",
			},
			providerUsername: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_USERNAME", ""),
				Description: "Couchbase username",
			},
			providerPassword: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_PASSWORD", ""),
				Sensitive:   true,
				Description: "Couchbase password",
			},
			providerConnectionTimeout: {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CB_MANAGEMENT_TIMEOUT", 15),
				Description: "Couchbase connection timeout",
			},
			providerTLSRootCertSkipVerify: {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TLS_ROOT_CERT_SKIP_VERIFY", false),
				Description: "TLS root certificate skip verify",
			},
			providerAllowSaslMechanism: {
				// We used TypeString with commas separator because TypeList doesn't support validate functionality
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TLS_ROOT_CERT_ALLOW_SASL_MECHANISM", "PLAIN,SCRAM-SHA1,SCRAM-SHA256,SCRAM-SHA512"),
				Description: fmt.Sprintf("Allowed Sasl Mechanism (separate values with commas \",\")\nAllowed values:\n%s\n%s\n%s\n%s\n",
					gocb.PlainSaslMechanism,
					gocb.ScramSha1SaslMechanism,
					gocb.ScramSha256SaslMechanism,
					gocb.ScramSha512SaslMechanism,
				),
				ValidateDiagFunc: validateAllowSaslMechanism(),
			},
			providerTLSRootCert: {
				Type:             schema.TypeString,
				Required:         false,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("TLS_ROOT_CERT", ""),
				Description:      "Path to TLS Root Certificate (in PEM format)",
				ValidateDiagFunc: validateTlsRootCert(),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"couchbase_bucket_manager":      resourceBucket(),
			"couchbase_security_group":      resourceSecurityGroup(),
			"couchbase_security_user":       resourceSecurityUser(),
			"couchbase_primary_query_index": resourcePrimaryQueryIndex(),
			"couchbase_query_index":         resourceQueryIndex(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

// getSaslMechanism function parse values from string (separeted with commas) and return
// list of gocb allow sasl mechanism
func getSaslMechanism(rawSaslMechanism string) ([]gocb.SaslMechanism, diag.Diagnostics) {
	var (
		saslMechanism []gocb.SaslMechanism
		diags         diag.Diagnostics
	)

	for _, value := range strings.Split(strings.ReplaceAll(rawSaslMechanism, " ", ""), ",") {
		switch gocb.SaslMechanism(value) {
		case gocb.PlainSaslMechanism:
			saslMechanism = append(saslMechanism, gocb.PlainSaslMechanism)
		case gocb.ScramSha1SaslMechanism:
			saslMechanism = append(saslMechanism, gocb.ScramSha1SaslMechanism)
		case gocb.ScramSha256SaslMechanism:
			saslMechanism = append(saslMechanism, gocb.ScramSha256SaslMechanism)
		case gocb.ScramSha512SaslMechanism:
			saslMechanism = append(saslMechanism, gocb.ScramSha512SaslMechanism)
		default:
			diags = append(diags, *getValidateAllowSaslMechanismDiagMessage(value))
			return nil, diags
		}
	}
	return saslMechanism, nil
}

// certificateManagement function add certificate from file to crypto/x509 certpool
func certificateManagement(filePath string) (*x509.CertPool, diag.Diagnostics) {
	tlsRootCAs := *x509.NewCertPool()

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	ok := tlsRootCAs.AppendCertsFromPEM(data)
	if !ok {
		return nil, diag.Errorf("cannot append certificate")
	}

	return &tlsRootCAs, nil
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var (
		tlsRootCAs *x509.CertPool
		diags      diag.Diagnostics
		scheme     string
	)

	certificatePath := d.Get(providerTLSRootCert).(string)
	if certificatePath != "" {
		tlsRootCAs, diags = certificateManagement(certificatePath)
		if diags != nil {
			return nil, diags
		}
		scheme = "couchbases"
	} else {
		scheme = "couchbase"
	}

	saslMechanism, diags := getSaslMechanism(d.Get(providerAllowSaslMechanism).(string))
	if diags != nil {
		return nil, diags
	}

	cc := &CouchbaseConnection{
		// Currently gocb v2 doesn't support some operations so we must use couchbase/couchbases and node-to-node
		// ports for connection. You can read more about ports here:
		// https://docs.couchbase.com/server/current/install/install-ports.html
		Scheme:     scheme,
		Address:    d.Get(providerAddress).(string),
		NodePort:   d.Get(providerNodePort).(int),
		ClientPort: d.Get(providerClientPort).(int),
		ClusterOptions: gocb.ClusterOptions{
			Username: d.Get(providerUsername).(string),
			Password: d.Get(providerPassword).(string),
			TimeoutsConfig: gocb.TimeoutsConfig{
				ManagementTimeout: time.Duration(d.Get(providerConnectionTimeout).(int)) * time.Second,
			},
			SecurityConfig: gocb.SecurityConfig{
				TLSSkipVerify:         d.Get(providerTLSRootCertSkipVerify).(bool),
				TLSRootCAs:            tlsRootCAs,
				AllowedSaslMechanisms: saslMechanism,
			},
		},
	}

	_, diags = cc.ConnectionValidate()
	if diags != nil {
		return nil, diags
	}

	return cc, diags
}
