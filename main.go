package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-rmon/rmon"
)

func main() {
	var debug bool
	var address string

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.StringVar(&address, "address", "provider", "this value is used in the TF_REATTACH_PROVIDERS environment variable during debugging")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: address,
		ProviderFunc: rmon.Provider,
	}
	plugin.Serve(opts)
}
