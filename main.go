package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/jaybocc2/terraform-provider-zip/zipfile"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zipfile.Provider,
	})
}
