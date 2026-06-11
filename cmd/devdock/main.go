package main

import (
	"fmt"
	"os"

	"devdock/internal/cli"
	"devdock/internal/errors"
	"devdock/internal/home"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	cli.Version = Version
	cli.Commit = Commit
	cli.BuildDate = BuildDate

	if err := home.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize DevDock home: %v\n", err)
		os.Exit(1)
	}

	if err := cli.Execute(); err != nil {
		errors.HandleError(err)
	}
}
