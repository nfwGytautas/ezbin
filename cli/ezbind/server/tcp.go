package server

import (
	"fmt"
	"log"
	"net"

	"github.com/nfwGytautas/ezbin/ezbin/connection"
	"github.com/nfwGytautas/ezbin/shared"
)

func startTcpServer(config *DaemonConfig) {
	// Start TCP server
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", config.Server.Port))
	if err != nil {
		log.Fatalf("error listening: %v", err)
		return
	}

	cwd, err := shared.CurrentDirectory()
	if err != nil {
		log.Fatalf("error getting current directory: %v", err)
		return
	}

	log.Printf("accepting connections on: %s", ln.Addr().String())
	log.Println("tcp server properties:")
	log.Printf("	+ port: %v", config.Server.Port)
	log.Printf("	+ framesize: %v", config.Server.FrameSize)
	log.Printf("current directory: %s", cwd)

	err = connection.ServeP2C(ln, connection.P2CServeParameters{
		ConnectionPrivateKey: config.Connection.Private,
		ServerIdentity:       config.Identifier,
		FrameSize:            config.Server.FrameSize,
		Protocol:             config.Peer.Protocol,
		PackageDir:           config.Storage.Location,
	})
	if err != nil {
		log.Fatalf("error serving: %v", err)
	}
}

func handleConnection(conn net.Conn) {
	// Close the connection when we're done
	defer conn.Close()

	// Read incoming data
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the incoming data
	fmt.Printf("Received: %s", buf)
}
