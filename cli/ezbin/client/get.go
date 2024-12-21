package ez_client

import (
	"fmt"
	"log"

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

	// Get package info
	pckInfo, err := conn.GetPackageInfo(pck)
	if err != nil {
		return err
	}

	log.Println(pckInfo)

	return nil
}
