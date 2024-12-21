package connection

import (
	"encoding/json"
	"net"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/shared"
)

// Connection between client and peer
type connectionC2P struct {
	conn   net.Conn
	params C2PConnectionParameters

	buffer []byte
}

// Open a tcp connection to a peer
func (c *connectionC2P) open() error {
	conn, err := net.Dial("tcp", c.params.Peer.Address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.buffer = make([]byte, HANDSHAKE_BUFFER_SIZE)

	// Handshake with the peer
	err = c.Send(requests.HandshakeRequest{
		UserIdentifier: c.params.UserIdentifier,
	})

	return nil
}

// Close the connection
func (c *connectionC2P) Close() {
	c.conn.Close()
}

// Send data to the peer
func (c *connectionC2P) Send(req requests.Request) error {
	header := []byte(req.Header())

	if len(header) > HEADER_SIZE_BYTES {
		return ErrHeaderTooLarge
	}

	err := shared.WriteSubRange(c.buffer, 0, header)
	if err != nil {
		return err
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = shared.WriteSubRange(c.buffer, HEADER_SIZE_BYTES, data)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(c.buffer[:HEADER_SIZE_BYTES+len(data)])
	if err != nil {
		return err
	}

	return nil
}
