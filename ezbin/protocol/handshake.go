package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

const pCONNECTION_KEY_SIZE = 2048

// ezbin `handshake` protocol implementation
type HandshakeProtocol struct {
	PrivateKey string `json:"private" yaml:"private"`
	PublicKey  string `json:"public" yaml:"public"`
}

// Encrypt data using the public key
func (hp *HandshakeProtocol) Encrypt(data []byte) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(hp.PublicKey)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), data)
}

// Decrypt data
func (hp *HandshakeProtocol) Decrypt(data []byte) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(hp.PrivateKey)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), data)
}

// Get the name of the protocol
func (p *HandshakeProtocol) Name() string {
	return PROTOCOL_HANDSHAKE
}

// Get a shareable key
func (p *HandshakeProtocol) GetShareableKey() string {
	return p.PublicKey
}

// Set the encryption key
func (p *HandshakeProtocol) SetEncryptionKey(key string) {
	p.PublicKey = key
}

// Get the version of the protocol
func (p *HandshakeProtocol) Version() string {
	return "0.1.0"
}

// Create a new handshake protocol
func NewHandshakeProtocol() (*HandshakeProtocol, error) {
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

	return &HandshakeProtocol{
		PrivateKey: base64.StdEncoding.EncodeToString(privateKeyBytes),
		PublicKey:  base64.StdEncoding.EncodeToString(publicKeyBytes),
	}, nil
}
