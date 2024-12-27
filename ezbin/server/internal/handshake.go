package server_internal

import (
	"log"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
)

func (c *serverP2CClient) handshake() error {
	req := requests.HandshakeRequest{}

	res := requests.HandshakeResponse{
		Okay:      true, // TODO: user authentication
		Framesize: c.config.Server.FrameSize,
		Protocol:  c.config.Peer.Protocol,
	}

	err := c.frame.ToJSON(&req)
	if err != nil {
		return err
	}

	log.Println("handshake request received:", req)

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
