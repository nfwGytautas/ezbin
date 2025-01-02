package ezbin_server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/google/uuid"
	"github.com/nfwGytautas/ezbin/ezbin"
	"github.com/nfwGytautas/ezbin/shared"
)

const pCONNECTION_KEY_SIZE = 2048

// Daemon config
type DaemonConfig struct {
	Version    string `yaml:"version"`
	Identifier string `yaml:"identifier"`

	Connection struct {
		Public  string `yaml:"public"`
		Private string `yaml:"private"`
	} `yaml:"connection"`

	Server struct {
		Port      int `yaml:"port"`
		FrameSize int `yaml:"framesize"`
	} `yaml:"server"`

	Storage struct {
		Location string `yaml:"location"`
	} `yaml:"storage"`

	Peer *struct {
		Protocol string `yaml:"protocol"`
	} `yaml:"peer"`
}

func NewDefaultDaemonConfig() (*DaemonConfig, error) {
	dc := DaemonConfig{}

	dc.Version = ezbin.VERSION

	// Generate identifier
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	dc.Identifier = uuid.String()

	// Generate connection key
	privateKey, err := rsa.GenerateKey(rand.Reader, pCONNECTION_KEY_SIZE)
	if err != nil {
		return nil, err
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	// Generate the public key
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	// Set key info
	dc.Connection.Public = string(publicKeyBytes)
	dc.Connection.Private = string(privateKeyBytes)

	// Other properties
	dc.Server.Port = 32000
	dc.Server.FrameSize = 1024

	dc.Storage.Location = "packages/"

	return &dc, nil
}

// Load the daemon config
func LoadDaemonConfig(config string) (*DaemonConfig, error) {
	dc := DaemonConfig{}

	err := shared.ReadYAML(config, &dc)
	if err != nil {
		return nil, err
	}

	return &dc, nil
}

// Save the daemon config
func (dc *DaemonConfig) Save() error {
	return shared.WriteYAML("ezbin.yaml", dc)
}

// Check if config is valid, returns true if the config is valid, false otherwise
func (dc *DaemonConfig) Validate() bool {
	return true
}

// Listen for incoming connections
func (dc *DaemonConfig) Run() error {
	log.Println("ezbin server starting...")
	log.Println("version: ", ezbin.VERSION)

	log.Printf("running %s", dc.Identifier)
	log.Printf("server connection key: %s", strings.ReplaceAll(dc.Connection.Public, "\n", ""))

	err := dc.initPackageDirectory()
	if err != nil {
		log.Fatalf("error initializing package directory: %v", err)
		return err
	}

	log.Printf("serving packages from: %s", dc.Storage.Location)

	// Start TCP server
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", dc.Server.Port))
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
	log.Printf("	+ port: %v", dc.Server.Port)
	log.Printf("	+ framesize: %v", dc.Server.FrameSize)
	log.Printf("current directory: %s", cwd)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleNewConnection(conn, dc)
	}
}

func (dc *DaemonConfig) initPackageDirectory() error {
	exists, err := shared.DirectoryExists(dc.Storage.Location)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	// Create directory
	err = shared.CreateDirectory(dc.Storage.Location)
	if err != nil {
		return err
	}

	// Create directory
	err = shared.CreateDirectory(dc.Storage.Location + "/.ezbin")
	if err != nil {
		return err
	}

	return nil
}
