package server

import (
	"log"
	"strings"
)

func RunServer(configFile string) {
	// Load config
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running %s", config.Identifier)
	log.Printf("server connection key: %s", strings.ReplaceAll(config.Connection.Public, "\n", ""))

	// Start TCP server
	startTcpServer(config)
}
