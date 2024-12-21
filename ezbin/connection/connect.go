package connection

// PeerConnectionData is a struct that represents a peer's connection data
type PeerConnectionData struct {
	// Address of the peer
	Address string

	// Connection key of the peer
	ConnectionKey string
}

// Arguments for connect function
type ConnectArgs struct {
	// Peer address
	Peer PeerConnectionData

	// User identifier
	UserIdentifier string
}

// Connect client to a peer
func ConnectC2P(args ConnectArgs) (C2PConnection, error) {
	// Connect to peer
	return C2PConnection{}, nil
}
