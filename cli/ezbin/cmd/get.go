package cmd

import (
	ez_client "github.com/nfwGytautas/ezbin/cli/ezbin/client"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <peer> <package>",
	Short: "Get a package",
	Long:  `Get a package`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ez_client.LoadUserIdentity()
		if err != nil {
			panic(err)
		}

		err = ez_client.GetPackage(identity, args[1], args[0])
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
