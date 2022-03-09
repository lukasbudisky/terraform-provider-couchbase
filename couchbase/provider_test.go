package couchbase

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"couchbase": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"CB_ADDRESS", "CB_CLIENT_PORT", "CB_NODE_PORT", "CB_USERNAME", "CB_PASSWORD"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("CB_ADDRESS, CB_CLIENT_PORT, CB_NODE_PORT, CB_USERNAME, CB_PASSWORD and optionally CB_MANAGEMENT_TIMEOUT must be set for acceptance tests")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(120)*time.Second)
	defer cancel()
	err := testAccProvider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}

}
