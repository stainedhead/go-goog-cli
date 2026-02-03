// Package main is the entry point for the goog CLI application.
package main

import (
	"os"

	"github.com/stainedhead/go-goog-cli/internal/adapter/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
