package main

import (
	"daas-terraform-provider-canal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: provider.Provider}

	plugin.Serve(opts)
}
