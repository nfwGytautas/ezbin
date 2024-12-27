package client

import (
	"fmt"
	"net"
	"os"

	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/errors"
	"github.com/nfwGytautas/ezbin/shared"
)

// PeerConnectionData is a struct that represents a peer's connection data
type PeerConnectionData struct {
	// Address of the peer
	Address string

	// Connection key of the peer
	ConnectionKey string
}

// Arguments for connect function
type C2PConnectionParameters struct {
	// Peer address
	Peer PeerConnectionData

	// User identifier
	UserIdentifier string
}

// Connection between client and peer
type ConnectionC2P struct {
	conn   net.Conn
	params C2PConnectionParameters

	frame     *connection.Frame
	framesize int
}

// Connect client to a peer
func Create(args C2PConnectionParameters) (*ConnectionC2P, error) {
	conn := ConnectionC2P{
		params: args,
	}

	err := conn.open()
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

// Open a tcp connection to a peer
func (c *ConnectionC2P) open() error {
	conn, err := net.Dial("tcp", c.params.Peer.Address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.frame = connection.NewFrame(conn, make([]byte, connection.HANDSHAKE_BUFFER_SIZE))

	// Handshake with the peer
	err = c.handshake()
	if err != nil {
		return err
	}

	return nil
}

// Close the connection
func (c *ConnectionC2P) Close() {
	c.conn.Close()
}

// Send data to the peer
func (c *ConnectionC2P) send(header string, req any) error {
	err := c.frame.SetHeader(header)
	if err != nil {
		return err
	}

	err = c.frame.FromJSON(req)
	if err != nil {
		return err
	}

	err = c.frame.Write()
	if err != nil {
		return err
	}

	return nil
}

// Start packet receive stream
func (c *ConnectionC2P) startPacketStream() error {
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

// Receive data from peer
func (c *ConnectionC2P) receive(res any) error {
	err := c.frame.Read()
	if err != nil {
		return err
	}

	header := c.frame.GetHeader()

	if header == requests.ERROR_HEADER {
		return errors.ErrIncorrectHeader
	}

	err = c.frame.ToJSON(res)
	if err != nil {
		return err
	}

	return nil
}

// Receive packet from peer
func (c *ConnectionC2P) receivePacket() error {
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

	return err
}

// Handshake with the peer
func (c *ConnectionC2P) handshake() error {
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

	c.frame = connection.NewFrame(c.conn, make([]byte, res.Framesize))
	c.framesize = res.Framesize

	return nil
}

// Get package info from peer
func (c *ConnectionC2P) GetPackageInfo(name string) (*requests.PackageInfoResponse, error) {
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

// Download package from peer
func (c *ConnectionC2P) DownloadPackage(name string, version string, outDir string, info *requests.PackageInfoResponse) error {
	req := requests.PackageDownloadRequest{
		Package: name,
		Version: version,
	}

	res := requests.PackageDownloadResponse{}

	err := c.send(requests.HeaderDownloadPackage, req)
	if err != nil {
		return err
	}

	err = c.receive(&res)
	if err != nil {
		return err
	}

	// Notify peer that we are ready to receive packets
	err = c.startPacketStream()
	if err != nil {
		return err
	}

	fmt.Println("Downloading package:", name)

	err = shared.CreateDirectory(outDir + ".ezbin")
	if err != nil {
		return err
	}

	// Create file
	filePath := outDir + ".ezbin/" + name + "@" + version + ".tar.gz"
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	totalReceivedSum := 0
	for i := 0; i < int(res.PacketCount); i++ {
		// TODO: Spanner
		err := c.frame.Read()
		if err != nil {
			// TODO: Error handling, request the same packet again
			return err
		}

		totalReceivedSum += c.frame.GetNumReadBytes()
		percentage := float64(totalReceivedSum) / float64(res.FullSize) * 100
		fmt.Printf("[%v|%v] Received %vB/%vB %v%%\n", i+1, res.PacketCount, totalReceivedSum, res.FullSize, percentage)

		if c.frame.GetNumReadBytes() == 0 {
			break
		}

		err = c.frame.TransferToWriter(file)
		if err != nil {
			return err
		}
	}

	file.Close()

	// Untar the package
	outPath := outDir + name + "/v" + version + "/"
	fmt.Println("Extracting into: ", outPath)
	err = shared.TarExtractDirectory(filePath, outPath)
	if err != nil {
		return err
	}

	return nil
}

// Upload package to peer
func (c *ConnectionC2P) UploadPackage(name string, version string, packageFile string) error {
	req := requests.PackageUploadRequest{
		Package: name,
		Version: version,
	}

	res := requests.PackageUploadResponse{}

	// Get size
	size, err := shared.FileSize(packageFile)
	if err != nil {
		return err
	}

	req.FullSize = uint64(size)

	// Calculate packet count
	sendableCount := int64(c.framesize - connection.HEADER_SIZE_BYTES)
	req.PacketCount = uint32(size / sendableCount)

	if size%sendableCount != 0 {
		req.PacketCount++
	}

	err = c.send(requests.HeaderUploadPackage, req)
	if err != nil {
		return err
	}

	err = c.receive(&res)
	if err != nil {
		return err
	}

	if !res.Okay {
		return errors.ErrUploadFailed
	}

	// Notify peer that we are ready to send packets
	err = c.startPacketStream()
	if err != nil {
		return err
	}

	// Wait for start
	err = c.frame.Read()
	if err != nil {
		return err
	}

	header := c.frame.GetHeader()
	if header != requests.HeaderPacket {
		return errors.ErrIncorrectHeader
	}

	// Open the file
	file, err := os.OpenFile(packageFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("Uploading package:", name)

	totalSentSum := 0
	for i := 0; i < int(req.PacketCount); i++ {
		err := c.frame.TransferFromReader(file, 0)
		if err != nil {
			return err
		}

		totalSentSum += c.frame.GetFrameSize()
		percentage := float64(totalSentSum) / float64(req.FullSize) * 100
		fmt.Printf("[%v|%v] Sent %vB/%vB %v%%\n", i+1, req.PacketCount, totalSentSum, req.FullSize, percentage)

		err = c.frame.Write()
		if err != nil {
			return err
		}
	}

	return nil
}
