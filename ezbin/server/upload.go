package ezbin_server

import (
	"fmt"
	"log"
	"os"

	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/shared"
)

func (c *serverP2CClient) uploadPackage() error {
	req := requests.PackageUploadRequest{}

	res := requests.PackageUploadResponse{}

	err := c.frame.ToJSON(&req)
	if err != nil {
		return err
	}

	packagePath := req.Package

	log.Println("package upload request received:", req)

	// TODO: Prohibit uploading any package starting with '.'

	if c.config.Storage.Location[0:2] == "./" {
		cwd, err := shared.CurrentDirectory()
		if err != nil {
			return err
		}

		packagePath = cwd + "/" + c.config.Storage.Location[2:] + "/" + packagePath
	} else {
		packagePath = c.config.Storage.Location + "/" + packagePath
	}

	packagePath = packagePath + "/v" + req.Version

	tempPath := c.config.Storage.Location + "/.ezbin/" // Temporary path

	tempPath = tempPath + req.Package + "@" + req.Version + ".tar.gz"

	// Create file
	file, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Send response
	res.Okay = true

	err = c.frame.FromJSON(res)
	if err != nil {
		return err
	}

	err = c.frame.Write()
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

	// Get data
	err = c.receivePacketStream()
	if err != nil {
		return err
	}

	totalReceivedSum := 0
	log.Println("receiving package packets...")
	for i := 0; i < int(req.PacketCount); i++ {
		err := c.receivePacket()
		if err != nil {
			// TODO: Error handling, request the same packet again
			return err
		}

		totalReceivedSum += c.frame.GetNumReadBytes()
		percentage := float64(totalReceivedSum) / float64(req.FullSize) * 100
		fmt.Printf("received packet: [%v|%v] %vB/%vB %v%%\n", i+1, req.PacketCount, totalReceivedSum, req.FullSize, percentage)

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
	err = shared.TarExtractDirectory(tempPath, packagePath)
	if err != nil {
		return err
	}

	return nil
}
