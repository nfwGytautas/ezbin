package ezbin_server

import (
	"log"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/protocol"
)

func (c *serverP2CClient) handshake() error {
	req := requests.HandshakeRequest{}

	res := requests.HandshakeResponse{
		Okay:      true, // TODO: user authentication
		Framesize: c.config.Server.FrameSize,
	}

	err := c.frame.Decrypt(c.config.Handshake.Decrypt)
	if err != nil {
		return err
	}

	err = c.frame.ToJSON(&req)
	if err != nil {
		return err
	}

	log.Println("handshake request received:", req)

	c.aesTransfer = protocol.NewAesTransferFromKey(req.Key)

	err = c.frame.FromJSON(res)
	if err != nil {
		return err
	}

	log.Println("handshake response sent:", res)

	err = c.frame.Encrypt(c.aesTransfer.Encrypt)
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
