package ezbin_server

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/protocol"
)

// TODO: Protocol

type serverP2CClient struct {
	config *DaemonConfig
	conn   net.Conn

	clientIdentity    string
	handshakeFinished bool
	frame             *connection.Frame
	aesTransfer       *protocol.AesTransfer
}

// Create new server client
func handleNewConnection(conn net.Conn, config *DaemonConfig) {
	client := &serverP2CClient{
		config:            config,
		conn:              conn,
		frame:             connection.NewFrame(conn, make([]byte, config.Server.FrameSize)),
		handshakeFinished: false,
	}

	client.handleConnection()
}

// Handle incoming connection
func (c *serverP2CClient) handleConnection() {
	defer c.conn.Close()

	log.Println("handling incoming connection...")

	for {
		err := c.conn.SetReadDeadline(time.Now().Add(connection.HANDSHAKE_READ_TIMEOUT))
		if err != nil {
			log.Fatal(err)
			return
		}

		err = c.frame.Read()
		if err != nil {
			if err == net.ErrClosed || err == io.EOF {
				log.Printf("connection closed with: %s", c.clientIdentity)
				return
			}

			log.Println("Error:", err)
			return
		}

		if !c.handshakeFinished {
			err := c.handshake()
			if err != nil {
				log.Println(err)
				return
			}

			continue
		}

		// Unmarshal the handshake request
		err = c.handleRequest()
		if err != nil {
			// TODO: Error handling
			log.Println(err)
			return
		}
	}
}

// Receive packet stream
func (c *serverP2CClient) receivePacketStream() error {
	err := c.frame.SetHeader(requests.HeaderPacket)
	if err != nil {
		return err
	}

	err = c.frame.Write()
	if err != nil {
		return err
	}

	return nil
}

// Receive packet
func (c *serverP2CClient) receivePacket() error {
	err := c.frame.Read()
	if err != nil {
		return err
	}

	header := c.frame.GetHeader()

	if header == requests.ERROR_HEADER {
		return ezbin.ErrIncorrectHeader
	}

	if header != requests.HeaderPacket {
		return ezbin.ErrIncorrectHeader
	}

	return nil
}

// Handle incoming request
func (c *serverP2CClient) handleRequest() error {
	err := c.frame.Decrypt(c.aesTransfer.Decrypt)
	if err != nil {
		return err
	}

	header := c.frame.GetHeader()

	log.Println("header received:", header)

	switch header {
	case requests.HeaderPackageInfo:
		return c.packageInfo()
	case requests.HeaderDownloadPackage:
		return c.downloadPackage()
	case requests.HeaderUploadPackage:
		return c.uploadPackage()
	}

	log.Println("unknown header received:", header)
	return ezbin.ErrUnknownHeader
}
