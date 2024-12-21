package cmd

import (
	"fmt"

	ez_client "github.com/nfwGytautas/ezbin/cli/ezbin/client"
	"github.com/spf13/cobra"
)

var peerCmd = &cobra.Command{
	Use:   "peer",
	Short: "Make changes to your known peers list",
	Long:  `Add, remove, or list known peers`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var peerAddCmd = &cobra.Command{
	Use:   "add <name> <address> <connection_key>",
	Short: "Add a new peer to your known peers list",
	Long:  `Add a new peer to your known peers list`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		verify, err := cmd.Flags().GetBool("verify")
		if err != nil {
			fmt.Println(err)
			return
		}

		identity, err := ez_client.LoadUserIdentity()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = identity.AddPeer(args[0], args[1], args[2], verify)
		if err != nil {
			cmd.Println("❌ Failed to add peer: ", err)
			return
		}

		cmd.Println("✅ Peer added")
	},
}

var peerRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a peer from your known peers list",
	Long:  `Remove a peer from your known peers list`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ez_client.LoadUserIdentity()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = identity.RemovePeer(args[0])
		if err != nil {
			cmd.Println("❌ Failed to remove peer: ", err)
			return
		}

		cmd.Println("✅ Peer removed")
	},
}

var peerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all known peers",
	Long:  `List all known peers`,
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ez_client.LoadUserIdentity()
		if err != nil {
			fmt.Println(err)
			return
		}

		identity.ListPeers()
	},
}

var peerCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the peer list",
	Long:  `Ping all peers and check which of them respond`,
	Run: func(cmd *cobra.Command, args []string) {
		identity, err := ez_client.LoadUserIdentity()
		if err != nil {
			fmt.Println(err)
			return
		}

		identity.CheckPeers()
	},
}

func init() {
	peerAddCmd.Flags().BoolP("verify", "v", false, "Verify connection to peer")

	peerCmd.AddCommand(peerAddCmd)
	peerCmd.AddCommand(peerRemoveCmd)
	peerCmd.AddCommand(peerListCmd)
	peerCmd.AddCommand(peerCheckCmd)
	rootCmd.AddCommand(peerCmd)
}
