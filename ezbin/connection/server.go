package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/errors"
	"github.com/nfwGytautas/ezbin/shared"
)

// TODO: Protocol

// Connection between peer and client
type serverP2C struct {
	ln     net.Listener
	params P2CServeParameters
}

type serverP2CClient struct {
	server *serverP2C
	conn   net.Conn

	clientIdentity    string
	buffer            []byte
	numReadBytes      int
	handshakeFinished bool
}

// Listen for incoming connections
func (c *serverP2C) handle() {
	for {
		conn, err := c.ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		client := serverP2CClient{
			server:            c,
			conn:              conn,
			buffer:            make([]byte, c.params.FrameSize),
			handshakeFinished: false,
		}

		go client.handleConnection(conn)
	}
}

// Handle incoming connection
func (c *serverP2CClient) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("handling incoming connection...")

	for {
		err := c.conn.SetReadDeadline(time.Now().Add(HANDSHAKE_READ_TIMEOUT))
		if err != nil {
			log.Fatal(err)
			return
		}

		n, err := conn.Read(c.buffer)
		if err != nil {
			if err == net.ErrClosed || err == io.EOF {
				log.Printf("connection closed with: %s", c.clientIdentity)
				return
			}

			log.Println("Error:", err)
			return
		}

		c.numReadBytes = n

		// Unmarshal the handshake request
		header := c.buffer[:HEADER_SIZE_BYTES]
		headerString := string(header)

		err = c.handleRequest(headerString)
		if err != nil {
			// TODO: Error handling
			log.Println(err)
			return
		}
	}
}

// Handle incoming request
func (c *serverP2CClient) handleRequest(header string) error {
	trimmed := strings.TrimRight(header, "\x00")

	if !c.handshakeFinished {
		if trimmed != requests.HeaderHandshake {
			return errors.ErrHandshakeNotFinished
		}

		return c.handshake()
	}

	switch trimmed {
	case requests.HeaderPackageInfo:
		return c.packageInfo()
	}

	log.Println("unknown header received:", header)
	return errors.ErrUnknownHeader
}

// Write response
func (c *serverP2CClient) writeResponse(res any) error {
	data, err := json.Marshal(res)
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

// Handshake with the client
func (c *serverP2CClient) handshake() error {
	req := requests.HandshakeRequest{}

	res := requests.HandshakeResponse{
		Okay:      true, // TODO: user authentication
		Framesize: c.server.params.FrameSize,
		Protocol:  c.server.params.Protocol,
	}

	err := json.Unmarshal(c.buffer[HEADER_SIZE_BYTES:c.numReadBytes], &req)
	if err != nil {
		return err
	}

	log.Println("handshake request received:", req)

	err = c.writeResponse(res)
	if err != nil {
		return err
	}

	c.clientIdentity = req.UserIdentifier
	c.handshakeFinished = true

	return nil
}

// Get information about a package
func (c *serverP2CClient) packageInfo() error {
	req := requests.PackageInfoRequest{}

	res := requests.PackageInfoResponse{}

	err := json.Unmarshal(c.buffer[HEADER_SIZE_BYTES:c.numReadBytes], &req)
	if err != nil {
		return err
	}

	packagePath := req.Package

	log.Println("package info request received:", req)

	if c.server.params.PackageDir[0:2] == "./" {
		cwd, err := shared.CurrentDirectory()
		if err != nil {
			return err
		}

		packagePath = cwd + "/" + c.server.params.PackageDir[2:] + "/" + packagePath
	} else {
		packagePath = c.server.params.PackageDir + "/" + packagePath
	}

	log.Println("checking package:", packagePath)

	// Check if package exists
	exists, err := shared.DirectoryExists(packagePath)
	if err != nil {
		return err
	}

	if exists {
		// Get package info
		size, err := shared.GetDirectorySize(packagePath)
		if err != nil {
			return err
		}

		log.Printf("package found: %s, size: %v", req.Package, size)
		res.Size = size
		res.Exists = true
	} else {
		log.Println("package not found:", req.Package)
		res.Exists = false
	}

	err = c.writeResponse(res)
	if err != nil {
		return err
	}

	return nil
}
