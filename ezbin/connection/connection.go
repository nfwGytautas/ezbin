package connection

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/errors"
	"github.com/nfwGytautas/ezbin/shared"
)

// Connection between client and peer
type connectionC2P struct {
	conn   net.Conn
	params C2PConnectionParameters

	buffer    []byte
	framesize int
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

// Write header to buffer
func (c *connectionC2P) writeHeader(header string) error {
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

// Get header
func (c *connectionC2P) getHeader() string {
	header := c.buffer[:HEADER_SIZE_BYTES]
	return strings.TrimRight(string(header), "\x00")
}

// Send data to the peer
func (c *connectionC2P) send(header string, req any) error {
	err := c.writeHeader(header)
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

// Start packet receive stream
func (c *connectionC2P) receivePacketStream() error {
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

// Start packet send stream
func (c *connectionC2P) sendPacketStream() error {
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

// Receive data from peer
func (c *connectionC2P) receive(res any) error {
	n, err := c.conn.Read(c.buffer)
	if err != nil {
		return err
	}

	header := c.getHeader()

	if header == requests.ERROR_HEADER {
		return errors.ErrIncorrectHeader
	}

	err = json.Unmarshal(c.buffer[HEADER_SIZE_BYTES:n], res)
	if err != nil {
		return err
	}

	return nil
}

// Receive packet from peer
func (c *connectionC2P) receivePacket() (int, error) {
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

	return n, err
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
	c.framesize = res.Framesize

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

// Download package from peer
func (c *connectionC2P) DownloadPackage(name string, version string, outDir string, info *requests.PackageInfoResponse) error {
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
	err = c.receivePacketStream()
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
		n, err := c.receivePacket()
		if err != nil {
			// TODO: Error handling, request the same packet again
			return err
		}

		totalReceivedSum += n - HEADER_SIZE_BYTES - PACKET_METADATA_SIZE
		percentage := float64(totalReceivedSum) / float64(res.FullSize) * 100
		fmt.Printf("[%v|%v] Received %vB/%vB %v%%\n", i+1, res.PacketCount, totalReceivedSum, res.FullSize, percentage)

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
	outPath := outDir + name + "/v" + version + "/"
	fmt.Println("Extracting into: ", outPath)
	err = shared.TarExtractDirectory(filePath, outPath)
	if err != nil {
		return err
	}

	return nil
}

// Upload package to peer
func (c *connectionC2P) UploadPackage(name string, version string, packageFile string) error {
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
	sendableCount := int64(c.framesize - HEADER_SIZE_BYTES - PACKET_METADATA_SIZE)
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
	err = c.sendPacketStream()
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
	file, err := os.OpenFile(packageFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("Uploading package:", name)

	totalSentSum := 0
	for i := 0; i < int(req.PacketCount); i++ {
		n, err := file.Read(c.buffer[HEADER_SIZE_BYTES+PACKET_METADATA_SIZE:])
		if err != nil {
			return err
		}

		totalSentSum += n
		percentage := float64(totalSentSum) / float64(req.FullSize) * 100
		fmt.Printf("[%v|%v] Sent %vB/%vB %v%%\n", i+1, req.PacketCount, totalSentSum, req.FullSize, percentage)

		_, err = c.conn.Write(c.buffer[:HEADER_SIZE_BYTES+PACKET_METADATA_SIZE+n])
		if err != nil {
			return err
		}
	}

	return nil
}
