package ezbin_client

import (
	"fmt"
	"net"
	"os"

	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/protocol"
	"github.com/nfwGytautas/ezbin/shared"
)

type client2P struct {
	conn        net.Conn
	frame       *connection.Frame
	framesize   int
	aesTransfer *protocol.AesTransfer
}

// Close the connection
func (c *client2P) Close() {
	c.conn.Close()
}

// Send data to the peer
func (c *client2P) send(header string, req any) error {
	err := c.frame.SetHeader(header)
	if err != nil {
		return err
	}

	err = c.frame.FromJSON(req)
	if err != nil {
		return err
	}

	err = c.frame.Encrypt(c.aesTransfer.Encrypt)
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
func (c *client2P) startPacketStream() error {
	err := c.frame.SetHeader(requests.HeaderPacket)
	if err != nil {
		return err
	}

	err = c.frame.Encrypt(c.aesTransfer.Encrypt)
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
func (c *client2P) receive(res any) error {
	err := c.frame.Read()
	if err != nil {
		return err
	}

	err = c.frame.Decrypt(c.aesTransfer.Decrypt)
	if err != nil {
		return err
	}

	header := c.frame.GetHeader()

	if header == requests.ERROR_HEADER {
		return ezbin.ErrIncorrectHeader
	}

	err = c.frame.ToJSON(res)
	if err != nil {
		return err
	}

	return nil
}

// Receive packet from peer
func (c *client2P) receivePacket() error {
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

	return err
}

// Handshake with the peer
func (c *client2P) handshake(identifier string, serverKey string, aesTransferKey string) error {
	req := requests.HandshakeRequest{
		UserIdentifier: identifier,
		Key:            aesTransferKey,
	}

	res := requests.HandshakeResponse{}

	hs := protocol.NewHandshakeFromKeys(serverKey, "")
	c.aesTransfer = protocol.NewAesTransferFromKey(aesTransferKey)

	// Send out handshake
	err := c.frame.SetHeader(requests.HeaderHandshake)
	if err != nil {
		return err
	}

	err = c.frame.FromJSON(req)
	if err != nil {
		return err
	}

	err = c.frame.Encrypt(hs.Encrypt)
	if err != nil {
		return err
	}

	err = c.frame.Write()
	if err != nil {
		return err
	}

	// Receive response
	err = c.receive(&res)
	if err != nil {
		return err
	}

	if res.Okay == false {
		return ezbin.ErrHandshakeFailed
	}

	c.frame = connection.NewFrame(c.conn, make([]byte, res.Framesize))
	c.framesize = res.Framesize

	return nil
}

// Get package info from peer
func (c *client2P) GetPackageInfo(name string) (*requests.PackageInfoResponse, error) {
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
func (c *client2P) DownloadPackage(name string, version string, outDir string, info *requests.PackageInfoResponse) error {
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

		err = c.frame.Decrypt(c.aesTransfer.Decrypt)
		if err != nil {
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
func (c *client2P) UploadPackage(name string, version string, packageFile string) error {
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
		return ezbin.ErrUploadFailed
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
		return ezbin.ErrIncorrectHeader
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
