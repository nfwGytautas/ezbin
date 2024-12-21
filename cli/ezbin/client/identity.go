package ez_client

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/ezbin/protocol"
	"github.com/nfwGytautas/ezbin/shared"
)

const IDENTITY_FILE = ".ezbin.identity.json"

// UserIdentity is a struct that represents a user's identity
type UserIdentity struct {
	Version      string                                   `json:"version"`
	Identifier   string                                   `json:"identifier"`
	ProtocolInfo map[string]protocol.ProtocolData         `json:"protocolData"`
	Peers        map[string]connection.PeerConnectionData `json:"Peers"`
}

// Load local user identity from `.ezbin.identity.json`
func LoadUserIdentity() (*UserIdentity, error) {
	// Check if file exists
	homeDir, err := shared.HomeDirectory()
	if err != nil {
		return nil, err
	}

	fullPath := homeDir + "/" + IDENTITY_FILE

	// Check if exists
	exists, err := shared.FileExists(fullPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrIdentityNotFound
	}

	// Load identity
	var identity UserIdentity

	err = shared.ReadJson(fullPath, &identity)
	if err != nil {
		return nil, err
	}

	// Check version
	if identity.Version != VERSION {
		fmt.Println("⚠️ Identity was generated with an older version of ezbin. Please upgrade your identity.")
	}

	return &identity, nil
}

// Generate new user identity
func GenerateUserIdentity() (*UserIdentity, error) {
	fmt.Println("🔨 Constructing new identity...")

	identity := UserIdentity{
		Version:      VERSION,
		ProtocolInfo: make(map[string]protocol.ProtocolData),
		Peers:        make(map[string]connection.PeerConnectionData),
	}

	// Generate identifier
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	identity.Identifier = uuid.String()

	// Generate protocols
	for _, protocol := range ezbin.GetSupportedProtocols() {
		fmt.Printf("	+ %s (%s)\n", protocol.Name(), protocol.Version())
		data, err := protocol.GenerateNew()
		if err != nil {
			return nil, err
		}

		identity.ProtocolInfo[protocol.Name()] = data
	}

	err = identity.Save()
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

// Save identity to file
func (us *UserIdentity) Save() error {
	homeDir, err := shared.HomeDirectory()
	if err != nil {
		return err
	}

	fullPath := homeDir + "/" + IDENTITY_FILE

	err = shared.WriteJson(fullPath, us)
	if err != nil {
		return err
	}

	return nil
}

// List peers from identity
func (us *UserIdentity) ListPeers() {
	if len(us.Peers) == 0 {
		fmt.Println("⚠️ No peers in identity, try adding one with `ezbin peer add`")
		return
	}

	fmt.Println("Peers:")
	for name, c := range us.Peers {
		fmt.Printf("	- %s %s\n", name, c.Address)
	}
}

// Check if peer exists
func (us *UserIdentity) KnowsPeer(name string) bool {
	_, ok := us.Peers[name]
	return ok
}

// Add peer to identity
func (us *UserIdentity) AddPeer(name string, addr string, key string, verify bool) error {
	// Check if peer already exists
	if us.KnowsPeer(name) {
		return ErrPeerExists
	}

	// Try and connect to the peer
	if verify {
		fmt.Println("Verifying connection to peer...")
		connection, err := connection.ConnectC2P(connection.C2PConnectionParameters{
			Peer: connection.PeerConnectionData{
				Address:       addr,
				ConnectionKey: key,
			},
			UserIdentifier: us.Identifier,
		})
		if err != nil {
			return err
		}
		defer connection.Close()
	}

	us.Peers[name] = connection.PeerConnectionData{
		Address:       addr,
		ConnectionKey: key,
	}

	err := us.Save()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

// Remove peer from identity
func (us *UserIdentity) RemovePeer(addr string) error {
	if _, ok := us.Peers[addr]; !ok {
		return ErrPeerNotFound
	}

	delete(us.Peers, addr)
	err := us.Save()
	if err != nil {
		return err
	}

	return nil
}

// Check connection to all known peers
func (us *UserIdentity) CheckPeers() {
	fmt.Println("Checking connections to peers...")

	for peer, c := range us.Peers {
		conn, err := connection.ConnectC2P(connection.C2PConnectionParameters{
			Peer:           c,
			UserIdentifier: us.Identifier,
		})
		if err != nil {
			fmt.Printf("	- %s: ❌ %s\n", peer, err)
			continue
		}

		conn.Close()

		fmt.Printf("	- %s: ✅\n", peer)
	}
}
