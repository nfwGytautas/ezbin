package server

import (
	"log"
	"strings"

	"github.com/nfwGytautas/ezbin/shared"
)

func RunServer(configFile string) {
	// Load config
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running %s", config.Identifier)
	log.Printf("server connection key: %s", strings.ReplaceAll(config.Connection.Public, "\n", ""))

	// Initialize package directory
	err = initPackageDirectory(config)

	log.Printf("serving packages from: %s", config.Storage.Location)

	// Start TCP server
	startTcpServer(config)
}

func initPackageDirectory(config *DaemonConfig) error {
	exists, err := shared.DirectoryExists(config.Storage.Location)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	// Create directory
	err = shared.CreateDirectory(config.Storage.Location)
	if err != nil {
		return err
	}

	// Create directory
	err = shared.CreateDirectory(config.Storage.Location + "/.ezbin")
	if err != nil {
		return err
	}

	return nil
}
