package ezbin

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/google/uuid"
	"github.com/nfwGytautas/ezbin/shared"
)

const VERSION = "0.1.0"
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

	dc.Version = VERSION

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
