package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lukasbudisky/terraform-provider-couchbase/couchbase"
)

func main() {
	const (
		hostname     = "budisky.com"
		namespace    = "couchbase"
		providerName = "couchbase"
	)

	var debugMode bool

	flag.BoolVar(
		&debugMode,
		"debug",
		false,
		"set to true to run the provider with support for debuggers like delve",
	)

	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: func() *schema.Provider {
		return couchbase.Provider()
	}}

	if debugMode {
		provider := fmt.Sprintf("%s/%s/%s", hostname, namespace, providerName)

		if err := plugin.Debug(
			context.Background(),
			provider,
			opts); err != nil {

			log.Fatal(err.Error())
		}
	}

	plugin.Serve(opts)
}
