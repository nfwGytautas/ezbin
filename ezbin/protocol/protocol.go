package protocol

import (
	"fmt"

	"github.com/nfwGytautas/ezbin/ezbin"
)

const (
	PROTOCOL_ECDSA     = "ecdsa"
	PROTOCOL_RSA       = "rsa"
	PROTOCOL_NOOP      = "no-op"
	PROTOCOL_HANDSHAKE = "ez-handshake"
)

// A struct containing all possible protocols
type Protocols struct {
	ECDSA *ECDSAProtocol `json:"ecdsa"`
	RSA   *RSAProtocol   `json:"rsa"`
	Noop  *NoopProtocol  `json:"noop"`
}

// Protocol interface
type Protocol interface {
	// Encrypt data
	Encrypt([]byte) ([]byte, error)

	// Decrypt data
	Decrypt([]byte) ([]byte, error)

	// Get name of the protocol
	Name() string

	// Get a shareable key
	GetShareableKey() string

	// Set encryption key
	SetEncryptionKey(string)
}

// Generate a protocols struct
func GenerateProtocols() (Protocols, error) {
	noop, err := NewNoOpProtocol()
	if err != nil {
		return Protocols{}, err
	}

	return Protocols{
		ECDSA: &ECDSAProtocol{},
		RSA:   &RSAProtocol{},
		Noop:  noop,
	}, nil
}

// String representation of the protocols
func (p Protocols) String() string {
	return fmt.Sprintf("Supported protocols:\n\tECDSA: %s\n\tRSA: %s\n\tNoop: %s", p.ECDSA, p.RSA, p.Noop)
}

// Get a protocol by name
func (p *Protocols) Get(name string) (Protocol, error) {
	switch name {
	// case PROTOCOL_ECDSA:
	// 	return &p.ECDSA, nil
	// case PROTOCOL_RSA:
	// 	return &p.RSA, nil
	case PROTOCOL_NOOP:
		return p.Noop, nil
	default:
		return nil, ezbin.ErrUnknownProtocol
	}
}
