package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lukasbudisky/terraform-provider-couchbase/couchbase"
)

func main() {
	var debugMode bool

	flag.BoolVar(
		&debugMode,
		"debug",
		false,
		"set to true to run the provider with support for debuggers like delve",
	)

	flag.Parse()

	// nolint:gocritic
	opts := &plugin.ServeOpts{ProviderFunc: func() *schema.Provider {
		return couchbase.Provider()
	}, Debug: debugMode}

	plugin.Serve(opts)

}
