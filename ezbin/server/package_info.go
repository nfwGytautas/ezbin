package ezbin_server

import (
	"log"

	"github.com/nfwGytautas/ezbin/ezbin/connection/requests"
	"github.com/nfwGytautas/ezbin/shared"
)

func (c *serverP2CClient) packageInfo() error {
	req := requests.PackageInfoRequest{}

	res := requests.PackageInfoResponse{}

	err := c.frame.ToJSON(&req)
	if err != nil {
		return err
	}

	packagePath := req.Package

	log.Println("package info request received:", req)

	if c.config.Storage.Location[0:2] == "./" {
		cwd, err := shared.CurrentDirectory()
		if err != nil {
			return err
		}

		packagePath = cwd + "/" + c.config.Storage.Location[2:] + "/" + packagePath
	} else {
		packagePath = c.config.Storage.Location + "/" + packagePath
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

	err = c.frame.FromJSON(res)
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
