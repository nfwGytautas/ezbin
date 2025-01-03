package ezbin_server

import (
	"log"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
)

func (c *serverP2CClient) handshake() error {
	req := requests.HandshakeRequest{}

	res := requests.HandshakeResponse{
		Okay:      true, // TODO: user authentication
		Framesize: c.config.Server.FrameSize,
	}

	err := c.frame.ToJSON(&req)
	if err != nil {
		return err
	}

	log.Println("handshake request received:", req)

	// Setup protocol
	p, err := c.config.Peer.Protocol.Get(req.Protocol)
	if err != nil {
		return err
	}

	// Set encryption key to the client key
	p.SetEncryptionKey(req.Key)
	c.frame.SetProtocol(p)

	err = c.frame.FromJSON(res)
	if err != nil {
		return err
	}

	err = c.frame.Write()
	if err != nil {
		return err
	}

	c.clientIdentity = req.UserIdentifier
	c.handshakeFinished = true

	return nil
}
