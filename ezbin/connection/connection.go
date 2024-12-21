package connection

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/errors"
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
	err = c.handshake()
	if err != nil {
		return err
	}

	return nil
}

// Close the connection
func (c *connectionC2P) Close() {
	c.conn.Close()
}

// Send data to the peer
func (c *connectionC2P) send(header string, req any) error {
	if len(header) > HEADER_SIZE_BYTES {
		return errors.ErrHeaderTooLarge
	}

	err := shared.WriteSubRange(c.buffer, 0, []byte(header))
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

// Receive data from peer
func (c *connectionC2P) receive(res any) error {
	n, err := c.conn.Read(c.buffer)
	if err != nil {
		return err
	}

	header := c.buffer[:HEADER_SIZE_BYTES]
	headerString := string(header)

	if strings.TrimRight(headerString, "\x00") == requests.ERROR_HEADER {
		return errors.ErrIncorrectHeader
	}

	err = json.Unmarshal(c.buffer[HEADER_SIZE_BYTES:n], res)
	if err != nil {
		return err
	}

	return nil
}

// Handshake with the peer
func (c *connectionC2P) handshake() error {
	req := requests.HandshakeRequest{
		UserIdentifier: c.params.UserIdentifier,
	}

	res := requests.HandshakeResponse{}

	err := c.send(requests.HeaderHandshake, req)
	if err != nil {
		return err
	}

	err = c.receive(&res)
	if err != nil {
		return err
	}

	if res.Okay == false {
		return errors.ErrHandshakeFailed
	}

	c.buffer = make([]byte, res.Framesize)

	return nil
}

// Get package info from peer
func (c *connectionC2P) GetPackageInfo(name string) (*requests.PackageInfoResponse, error) {
	req := requests.PackageInfoRequest{
		Package: name,
	}

	res := requests.PackageInfoResponse{}

	err := c.send(requests.HeaderPackageInfo, req)
	if err != nil {
		return nil, err
	}

	err = c.receive(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
