package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/litebitus/terraform-provider-dql/internal/dqlprovider"
)

var version = "dev"

func main() {
	err := providerserver.Serve(context.Background(), dqlprovider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/litebitus/dql",
	})
	if err != nil {
		log.Fatal(err)
	}
}
