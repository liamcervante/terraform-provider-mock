package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	fmt.Printf("Looooook at meeeee")
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "terraform.local/local/fakelocal",
		Debug:   debug,
	}

	ctx := context.Background()
	tflog.Trace(ctx, "Running main function of provider")

	err := providerserver.Serve(ctx, provider.New(version), opts)

	if err != nil {
		fmt.Printf(err.Error())
		log.Fatal(err.Error())
	}
}