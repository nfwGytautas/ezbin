package server_internal

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/errors"
)

// TODO: Protocol

type serverP2CClient struct {
	config *ezbin.DaemonConfig
	conn   net.Conn

	clientIdentity    string
	handshakeFinished bool
	frame             *connection.Frame
}

// Create new server client
func Handle(conn net.Conn, config *ezbin.DaemonConfig) {
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
		return errors.ErrIncorrectHeader
	}

	if header != requests.HeaderPacket {
		return errors.ErrIncorrectHeader
	}

	return nil
}

// Handle incoming request
func (c *serverP2CClient) handleRequest() error {
	header := c.frame.GetHeader()

	if !c.handshakeFinished {
		if header != requests.HeaderHandshake {
			return errors.ErrHandshakeNotFinished
		}

		return c.handshake()
	}

	switch header {
	case requests.HeaderPackageInfo:
		return c.packageInfo()
	case requests.HeaderDownloadPackage:
		return c.downloadPackage()
	case requests.HeaderUploadPackage:
		return c.uploadPackage()
	}

	log.Println("unknown header received:", header)
	return errors.ErrUnknownHeader
}
