package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
		err = c.handleRequest()
		if err != nil {
			// TODO: Error handling
			log.Println(err)
			return
		}
	}
}

// Get header from buffer
func (c *serverP2CClient) getHeader() string {
	header := c.buffer[:HEADER_SIZE_BYTES]
	return strings.TrimRight(string(header), "\x00")
}

// Write header to buffer
func (c *serverP2CClient) writeHeader(header string) error {
	if len(header) > HEADER_SIZE_BYTES {
		return errors.ErrHeaderTooLarge
	}

	// Zero the header range
	err := shared.WriteSubRange(c.buffer, 0, []byte(strings.Repeat("\x00", HEADER_SIZE_BYTES)))
	if err != nil {
		return err
	}

	err = shared.WriteSubRange(c.buffer, 0, []byte(header))
	if err != nil {
		return err
	}

	return nil
}

// Receive packet stream
func (c *serverP2CClient) receivePacketStream() error {
	err := c.writeHeader(requests.HeaderPacket)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(c.buffer[:HEADER_SIZE_BYTES])
	if err != nil {
		return err
	}

	return nil
}

// Receive packet
func (c *serverP2CClient) receivePacket() (int, error) {
	n, err := c.conn.Read(c.buffer)
	if err != nil {
		return 0, err
	}

	header := c.getHeader()

	if header == requests.ERROR_HEADER {
		return 0, errors.ErrIncorrectHeader
	}

	if header != requests.HeaderPacket {
		return 0, errors.ErrIncorrectHeader
	}

	return n, nil
}

// Handle incoming request
func (c *serverP2CClient) handleRequest() error {
	header := c.getHeader()

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

// Download a package
func (c *serverP2CClient) downloadPackage() error {
	req := requests.PackageDownloadRequest{}

	res := requests.PackageDownloadResponse{}

	err := json.Unmarshal(c.buffer[HEADER_SIZE_BYTES:c.numReadBytes], &req)
	if err != nil {
		return err
	}

	packagePath := req.Package

	log.Println("package download request received:", req)

	// TODO: Prohibit getting any package starting with '.'

	if c.server.params.PackageDir[0:2] == "./" {
		cwd, err := shared.CurrentDirectory()
		if err != nil {
			return err
		}

		packagePath = cwd + "/" + c.server.params.PackageDir[2:] + "/" + packagePath
	} else {
		packagePath = c.server.params.PackageDir + "/" + packagePath
	}

	packagePath = packagePath + "/v" + req.Version

	tempPath := c.server.params.PackageDir + "/.ezbin/" // Temporary path

	// Prepare package by tarring/zipping it
	tempPath = tempPath + req.Package + "@" + req.Version + ".tar.gz"
	err = shared.TarCompressDirectory(packagePath, tempPath)
	if err != nil {
		return err
	}

	// Get package info
	size, err := shared.FileSize(tempPath)
	if err != nil {
		return err
	}

	sendableCount := int64(c.server.params.FrameSize - HEADER_SIZE_BYTES - PACKET_METADATA_SIZE)

	res.Okay = true
	res.PacketCount = uint32(size / sendableCount)
	res.FullSize = uint64(size)

	if size%sendableCount != 0 {
		res.PacketCount++
	}

	err = c.writeResponse(res)
	if err != nil {
		return err
	}

	// Wait for start
	_, err = c.conn.Read(c.buffer)
	if err != nil {
		return err
	}

	header := c.getHeader()
	if header != requests.HeaderPacket {
		return errors.ErrIncorrectHeader
	}

	// Open the file
	file, err := os.OpenFile(tempPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Start packet stream
	log.Printf("sending %v packets", res.PacketCount)
	for i := 0; i < int(res.PacketCount); i++ {
		n, err := file.Read(c.buffer[HEADER_SIZE_BYTES+PACKET_METADATA_SIZE:])
		if err != nil {
			return err
		}

		log.Printf("sending packet: [%v/%v] size (no header): %v for %s", i+1, res.PacketCount, n, c.clientIdentity)
		_, err = c.conn.Write(c.buffer[:HEADER_SIZE_BYTES+PACKET_METADATA_SIZE+n])
		if err != nil {
			return err
		}
	}

	return nil
}

// Upload a package
func (c *serverP2CClient) uploadPackage() error {
	req := requests.PackageUploadRequest{}

	res := requests.PackageUploadResponse{}

	err := json.Unmarshal(c.buffer[HEADER_SIZE_BYTES:c.numReadBytes], &req)
	if err != nil {
		return err
	}

	packagePath := req.Package

	log.Println("package upload request received:", req)

	// TODO: Prohibit uploading any package starting with '.'

	if c.server.params.PackageDir[0:2] == "./" {
		cwd, err := shared.CurrentDirectory()
		if err != nil {
			return err
		}

		packagePath = cwd + "/" + c.server.params.PackageDir[2:] + "/" + packagePath
	} else {
		packagePath = c.server.params.PackageDir + "/" + packagePath
	}

	packagePath = packagePath + "/v" + req.Version

	tempPath := c.server.params.PackageDir + "/.ezbin/" // Temporary path

	tempPath = tempPath + req.Package + "@" + req.Version + ".tar.gz"

	// Create file
	file, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Send response
	res.Okay = true

	err = c.writeResponse(res)
	if err != nil {
		return err
	}

	// Wait for start
	_, err = c.conn.Read(c.buffer)
	if err != nil {
		return err
	}

	header := c.getHeader()
	if header != requests.HeaderPacket {
		return errors.ErrIncorrectHeader
	}

	// Get data
	err = c.receivePacketStream()
	if err != nil {
		return err
	}

	totalReceivedSum := 0
	log.Println("receiving package packets...")
	for i := 0; i < int(req.PacketCount); i++ {
		n, err := c.receivePacket()
		if err != nil {
			// TODO: Error handling, request the same packet again
			return err
		}

		totalReceivedSum += n - HEADER_SIZE_BYTES - PACKET_METADATA_SIZE
		percentage := float64(totalReceivedSum) / float64(req.FullSize) * 100
		fmt.Printf("received packet: [%v|%v] %vB/%vB %v%%\n", i+1, req.PacketCount, totalReceivedSum, req.FullSize, percentage)

		if n == 0 {
			break
		}

		_, err = file.Write(c.buffer[HEADER_SIZE_BYTES+PACKET_METADATA_SIZE : n])
		if err != nil {
			return err
		}
	}

	file.Close()

	// Untar the package
	err = shared.TarExtractDirectory(tempPath, packagePath)
	if err != nil {
		return err
	}

	return nil
}
