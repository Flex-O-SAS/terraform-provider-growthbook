// Description: Main entry point for the Terraform provider plugin.
package main

import (
	"context"
	"log"

	"terraform-provider-growthbook/internal"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), internal.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/growthbook",
	})
	if err != nil {
		log.Fatal(err)
	}
}
