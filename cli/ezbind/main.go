package main

import (
	"log"
	"os"

	"github.com/nfwGytautas/ezbin/cli/ezbind/server"
	"github.com/nfwGytautas/ezbin/ezbin"
	ezbin_server "github.com/nfwGytautas/ezbin/ezbin/server"
	"github.com/nfwGytautas/ezbin/shared"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		log.Println("ezbin version: ", ezbin.VERSION)
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "generate-peer") {
		log.Println("Generating default peer config...")
		config, err := server.NewPeerConfig()
		if err != nil {
			log.Fatal(err)
		}

		err = config.Save()
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	configPath := "ezbin.yaml"

	if len(os.Args) > 1 {
		log.Println("using config: ", os.Args[1])
		configPath = os.Args[1]
	}

	// Check if `ezbin.yaml` exists
	exists, err := shared.FileExists(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	if !exists {
		log.Fatalf("%s not found", configPath)
		return
	}

	cfg, err := ezbin.LoadDaemonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = ezbin_server.New(*cfg).Run()
	if err != nil {
		log.Fatal(err)
	}
}
