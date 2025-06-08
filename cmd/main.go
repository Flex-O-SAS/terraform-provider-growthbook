// Description: Main entry point for the Terraform provider plugin.
package main

import (
	"terraform-provider-growthbook/internal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: internal.Provider,
	})
}
