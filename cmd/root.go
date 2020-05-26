package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	Version string
)

var rootCmd = &cobra.Command{
	Use:     "go-pivnet-client",
	Short:   "A client implementation of the Pivnet API for use of Data R&D to release GPDB on network.pivotal.io",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
