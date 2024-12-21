package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/google/uuid"
	"github.com/nfwGytautas/ezbin/shared"
)

const VERSION = "0.1.0"
const CONNECTION_KEY_SIZE = 2048

// Daemon config
type DaemonConfig struct {
	Version              string      `yaml:"version"`
	Identifier           string      `yaml:"identifier"`
	ConnectionKey        string      `yaml:"connectionKey"`
	ConnectionPrivateKey string      `yaml:"connectionPrivateKey"`
	Peer                 *PeerConfig `yaml:"peer"`
}

// Peer mode
type PeerConfig struct {
}

// Create a new base config
func newBaseConfig() (*DaemonConfig, error) {
	dc := DaemonConfig{}

	dc.Version = VERSION

	// Generate identifier
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	dc.Identifier = uuid.String()

	// Generate connection key
	privateKey, err := rsa.GenerateKey(rand.Reader, CONNECTION_KEY_SIZE)
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
	dc.ConnectionKey = string(publicKeyBytes)
	dc.ConnectionPrivateKey = string(privateKeyBytes)

	return &dc, nil
}

// Create a new default peer config
func NewPeerConfig() (*DaemonConfig, error) {
	dc, err := newBaseConfig()
	if err != nil {
		return nil, err
	}

	// Add peer config

	return dc, nil
}

// Load the config from file
func LoadConfig(path string) (*DaemonConfig, error) {
	dc := DaemonConfig{}

	err := shared.ReadYAML(path, &dc)
	if err != nil {
		return nil, err
	}

	return &dc, nil
}

// Save the daemon config
func (dc *DaemonConfig) Save() error {
	return shared.WriteYAML("ezbin.yaml", dc)
}
