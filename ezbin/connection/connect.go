package connection

import "net"

// PeerConnectionData is a struct that represents a peer's connection data
type PeerConnectionData struct {
	// Address of the peer
	Address string

	// Connection key of the peer
	ConnectionKey string
}

// Arguments for connect function
type C2PConnectionParameters struct {
	// Peer address
	Peer PeerConnectionData

	// User identifier
	UserIdentifier string
}

// Arguments for serve function
type P2CServeParameters struct {
	// Key used to decrypt initial handshake messages
	ConnectionPrivateKey string

	// Server identity
	ServerIdentity string

	// Frame size
	FrameSize int
}

// Connect client to a peer
func ConnectC2P(args C2PConnectionParameters) (*connectionC2P, error) {
	conn := connectionC2P{
		params: args,
	}

	err := conn.open()
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

// Connect a peer to a client
func ServeP2C(ln net.Listener, args P2CServeParameters) error {
	conn := serverP2C{
		ln:     ln,
		params: args,
	}

	conn.handle()

	return nil
}
