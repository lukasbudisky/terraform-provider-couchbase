package couchbase

import (
	"context"
	"os"
	"testing"

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
	for _, name := range []string{"CB_ADDRESS", "CB_PORT", "CB_USERNAME", "CB_PASSWORD"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("CB_ADDRESS, CB_PORT, CB_USERNAME, CB_PASSWORD and optionally CB_MANAGEMENT_TIMEOUT must be set for acceptance tests")
		}
	}

	// TODO
	err := testAccProvider.Configure(context.TODO(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}

}
