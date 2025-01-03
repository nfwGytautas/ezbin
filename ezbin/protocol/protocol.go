package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
)

const pCONNECTION_KEY_SIZE = 2048
const cAES_KEY_SIZE = 32

// Handshake protocol using RSA encryption
type Handshake struct {
	// Encryption key when sending data
	EncryptionKey string `json:"encryptionKey" yaml:"encryptionKey"`

	// Decryption key when receiving data
	DecryptionKey string `json:"decryptionKey" yaml:"decryptionKey"`
}

// AES transfer protocol
type AesTransfer struct {
	// Encryption key when sending data, this will be the remote key
	Key string `json:"key" yaml:"key"`
}

// Encrypt using the handshake protocol
func (h *Handshake) Encrypt(data []byte) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(h.EncryptionKey)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), data)
}

// Decrypt using the handshake protocol
func (h *Handshake) Decrypt(data []byte) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(h.DecryptionKey)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), data)
}

// Encrypt using the AES transfer protocol
func (a *AesTransfer) Encrypt(data []byte) ([]byte, error) {
	key, err := hex.DecodeString(a.Key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)
	return encryptedData, nil
}

// Decrypt using the AES transfer protocol
func (a *AesTransfer) Decrypt(data []byte) ([]byte, error) {
	key, err := hex.DecodeString(a.Key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := data[:gcm.NonceSize()]
	encryptedData := data[gcm.NonceSize():]

	decryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

// Create new handshake protocol
func NewHandshake() (*Handshake, error) {
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

	return &Handshake{
		EncryptionKey: base64.StdEncoding.EncodeToString(publicKeyBytes),
		DecryptionKey: base64.StdEncoding.EncodeToString(privateKeyBytes),
	}, nil
}

// Create new handshake protocol from keys (base64 encoded)
func NewHandshakeFromKeys(encryptionKey string, decryptionKey string) *Handshake {
	return &Handshake{
		EncryptionKey: encryptionKey,
		DecryptionKey: decryptionKey,
	}
}

// Create new AES transfer protocol
func NewAesTransfer() (*AesTransfer, error) {
	key := make([]byte, cAES_KEY_SIZE)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return &AesTransfer{
		Key: hex.EncodeToString(key),
	}, nil
}

// Create new AES transfer protocol from key (hex encoded)
func NewAesTransferFromKey(key string) *AesTransfer {
	return &AesTransfer{
		Key: key,
	}
}
