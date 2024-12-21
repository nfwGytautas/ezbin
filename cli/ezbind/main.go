package main

import (
	"log"
	"os"

	"github.com/nfwGytautas/ezbin/shared"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		log.Println("ezbin version: ", VERSION)
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "generate-peer") {
		log.Println("Generating default peer config...")
		config, err := NewPeerConfig()
		if err != nil {
			log.Fatal(err)
		}

		err = config.Save()
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	log.Println("ezbin daemon starting...")
	log.Println("version: ", VERSION)

	configPath := "ezbin.yaml"

	if len(os.Args) > 1 {
		log.Println("using config: ", os.Args[1])
		configPath = os.Args[1]
	}

	// Check if `ezbin.yaml` exists
	exists, err := shared.FileExists(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Fatalf("%s not found", configPath)
		return
	}

	runDaemon(configPath)
}

func runDaemon(configPath string) {
	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running %s", config.Identifier)

	// Start the daemon
}
