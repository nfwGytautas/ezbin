package cmd

import (
	ezbin_client "github.com/nfwGytautas/ezbin/ezbin/client"
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:       "package <command>",
	Short:     "Make changes to your package list",
	Long:      `Get, remove, or list known packages`,
	ValidArgs: []string{"get", "remove", "pub", "list"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}

		err := cmd.ValidateArgs(args)
		if err != nil {
			panic(err)
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get <peer> <package>",
	Short: "Get a package",
	Long:  `Get a package`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ezbin_client.LoadUserIdentity()
		if err != nil {
			panic(err)
		}

		err = identity.GetPackage(args[1], args[0])
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
		identity, err := ezbin_client.LoadUserIdentity()
		if err != nil {
			panic(err)
		}

		err = identity.RemovePackage(args[0])
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
		identity, err := ezbin_client.LoadUserIdentity()
		if err != nil {
			panic(err)
		}

		err = identity.ListPackages()
		if err != nil {
			panic(err)
		}
	},
}

var pubCmd = &cobra.Command{
	Use:   "pub <directory|file> <version> <peer>",
	Short: "Publish a package",
	Long:  `Publish a package`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ezbin_client.LoadUserIdentity()
		if err != nil {
			panic(err)
		}

		err = identity.PublishPackage(args[0], args[1], args[2])
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	packageCmd.AddCommand(getCmd)
	packageCmd.AddCommand(removeCmd)
	packageCmd.AddCommand(pubCmd)
	packageCmd.AddCommand(listCmd)
	rootCmd.AddCommand(packageCmd)
}
