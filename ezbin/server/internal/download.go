package server_internal

import (
	"log"
	"os"

	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/ezbin/errors"
	"github.com/nfwGytautas/ezbin/shared"
)

func (c *serverP2CClient) downloadPackage() error {
	req := requests.PackageDownloadRequest{}

	res := requests.PackageDownloadResponse{}

	err := c.frame.ToJSON(&req)
	if err != nil {
		return err
	}

	packagePath := req.Package

	log.Println("package download request received:", req)

	// TODO: Prohibit getting any package starting with '.'

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

	log.Printf("package compressed size: %vB\n", size)

	sendableCount := int64(c.config.Server.FrameSize - connection.HEADER_SIZE_BYTES)

	res.Okay = true
	res.PacketCount = uint32(size / sendableCount)
	res.FullSize = uint64(size)

	if size%sendableCount != 0 {
		res.PacketCount++
	}

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
		err := c.frame.TransferFromReader(file, 0)
		if err != nil {
			return err
		}

		log.Printf("sending packet: [%v/%v] size (no header): %v for %s", i+1, res.PacketCount, c.frame.GetFrameSize(), c.clientIdentity)
		err = c.frame.Write()
		if err != nil {
			return err
		}
	}

	return nil
}
