package main

import (
	"os"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
)

func main() {
	viewKitCli := cli.NewViewCreatorCommand()
	if err := viewKitCli.Execute(); err != nil {
		// this error is okay to discard because cobra
		// logs any errors encountered during execution
		//
		// exiting with a non-zero status code signals
		// that an error has ocurred during execution
		os.Exit(1)
	}
}
