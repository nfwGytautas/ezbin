package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
)

// Connection between peer and client
type serverP2C struct {
	ln     net.Listener
	params P2CServeParameters
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
		go c.handleConnection(conn)
	}
}

// Handle incoming connection
func (c *serverP2C) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("handling incoming connection...")

	err := c.handshake(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a buffer to read data into
	buffer := make([]byte, c.params.FrameSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Received: %s\n", buffer[:n])
	}
}

// Handshake with the client
func (c *serverP2C) handshake(conn net.Conn) error {
	buffer := make([]byte, HANDSHAKE_BUFFER_SIZE)

	// Wait for the client to send the handshake request
	err := conn.SetReadDeadline(time.Now().Add(HANDSHAKE_READ_TIMEOUT))
	if err != nil {
		return err
	}

	// Read the handshake request
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	// Unmarshal the handshake request
	header := buffer[:HEADER_SIZE_BYTES]
	headerString := string(header)

	req, _ := requests.HeaderToRequestResponse(headerString)
	if req == nil {
		return ErrUnknownHeader
	}

	err = json.Unmarshal(buffer[HEADER_SIZE_BYTES:n], &req)
	if err != nil {
		return err
	}

	log.Println("Handshake request received:", req)

	// Remove deadline
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		return err
	}

	return nil
}
