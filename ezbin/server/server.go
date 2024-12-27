package ezbin_server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/nfwGytautas/ezbin/ezbin"
	server_internal "github.com/nfwGytautas/ezbin/ezbin/server/internal"
	"github.com/nfwGytautas/ezbin/shared"
)

// Connection between peer and client
type serverP2C struct {
	config ezbin.DaemonConfig
}

// Create new server from config
func New(config ezbin.DaemonConfig) *serverP2C {
	return &serverP2C{
		config: config,
	}
}

// Listen for incoming connections
func (c *serverP2C) Run() error {
	log.Println("ezbin server starting...")
	log.Println("version: ", ezbin.VERSION)

	log.Printf("running %s", c.config.Identifier)
	log.Printf("server connection key: %s", strings.ReplaceAll(c.config.Connection.Public, "\n", ""))

	err := c.initPackageDirectory()
	if err != nil {
		log.Fatalf("error initializing package directory: %v", err)
		return err
	}

	log.Printf("serving packages from: %s", c.config.Storage.Location)

	// Start TCP server
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", c.config.Server.Port))
	if err != nil {
		log.Fatalf("error listening: %v", err)
		return err
	}

	cwd, err := shared.CurrentDirectory()
	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return err
	}

	log.Printf("accepting connections on: %s", ln.Addr().String())
	log.Println("tcp server properties:")
	log.Printf("	+ port: %v", c.config.Server.Port)
	log.Printf("	+ framesize: %v", c.config.Server.FrameSize)
	log.Printf("current directory: %s", cwd)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go server_internal.Handle(conn, &c.config)
	}
}

func (s *serverP2C) initPackageDirectory() error {
	exists, err := shared.DirectoryExists(s.config.Storage.Location)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	// Create directory
	err = shared.CreateDirectory(s.config.Storage.Location)
	if err != nil {
		return err
	}

	// Create directory
	err = shared.CreateDirectory(s.config.Storage.Location + "/.ezbin")
	if err != nil {
		return err
	}

	return nil
}
