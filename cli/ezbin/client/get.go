package ez_client

import (
	"fmt"

	"github.com/nfwGytautas/ezbin/ezbin/connection"
)

func GetPackage(i *UserIdentity, pck string, peer string) error {
	fmt.Printf("📦 Getting package: %v\n", pck)

	if !i.KnowsPeer(peer) {
		return ErrPeerNotFound
	}

	// Open a connection to peer
	connData := i.Peers[peer]

	conn, err := connection.ConnectC2P(connection.C2PConnectionParameters{
		Peer:           connData,
		UserIdentifier: i.Identifier,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
