package cmd

import (
	"os"

	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ezbin",
	Short:   "ezbin CLI client",
	Long:    "ezbin is a CLI client for the ezbin artifactory service",
	Version: ezbin.VERSION,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
