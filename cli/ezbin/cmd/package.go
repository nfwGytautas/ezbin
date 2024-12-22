package cmd

import (
	ez_client "github.com/nfwGytautas/ezbin/cli/ezbin/client"
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Make changes to received packages",
	Long:  `Get, remove, or list known packages`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

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

var removeCmd = &cobra.Command{
	Use:   "remove <package>",
	Short: "Remove a package",
	Long:  `Remove a package`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ez_client.LoadUserIdentity()
		if err != nil {
			panic(err)
		}

		err = ez_client.RemovePackage(identity, args[0])
		if err != nil {
			panic(err)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all known packages",
	Long:  `List all known packages`,
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ez_client.LoadUserIdentity()
		if err != nil {
			panic(err)
		}

		err = ez_client.ListPackages(identity)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	packageCmd.AddCommand(getCmd)
	packageCmd.AddCommand(removeCmd)
	packageCmd.AddCommand(listCmd)
	rootCmd.AddCommand(packageCmd)
}
