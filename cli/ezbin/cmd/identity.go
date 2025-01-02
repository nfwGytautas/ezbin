package cmd

import (
	"errors"

	"github.com/nfwGytautas/ezbin/ezbin"
	ezbin_client "github.com/nfwGytautas/ezbin/ezbin/client"
	"github.com/spf13/cobra"
)

var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Make changes to your identity",
	Long:  `Create, check, update, migrate, import an identity`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var checkIdentity = &cobra.Command{
	Use:   "check",
	Short: "Check your identity",
	Long:  `Check your identity`,
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ezbin_client.LoadUserIdentity()
		if err != nil {
			if errors.Is(err, ezbin.ErrIdentityNotFound) {
				cmd.Println("❌ Identity not found")
			} else {
				cmd.Println("❌ Error loading identity: ", err)
			}
			return
		}

		cmd.Println("✅ Identity loaded successfully")
		cmd.Println(identity.Identifier)
	},
}

var generateIdentity = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new identity",
	Long:  `Generate a new identity`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we already have a user identity
		_, err := ezbin_client.LoadUserIdentity()
		if err == nil {
			confirmed, err := promptConfirmation("❗️ You already have an identity. Are you sure you want to generate a new one?")
			if err != nil || !confirmed {
				cmd.Println("❌ Aborted")
				return
			}
		}

		_, err = ezbin_client.GenerateUserIdentity()
		if err != nil {
			cmd.Println("❌ Error generating identity: ", err)
			return
		}

		cmd.Println("✅ Identity generated successfully")
	},
}

func init() {
	identityCmd.AddCommand(checkIdentity)
	identityCmd.AddCommand(generateIdentity)
	rootCmd.AddCommand(identityCmd)
}
