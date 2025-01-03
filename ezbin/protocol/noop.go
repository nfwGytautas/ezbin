package protocol

// Noop protocol implementation
type NoopProtocol struct {
}

// Encrypt data
func (p *NoopProtocol) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

// Decrypt data
func (p *NoopProtocol) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}

// Get the name of the protocol
func (p *NoopProtocol) Name() string {
	return PROTOCOL_NOOP
}

// Get a shareable key
func (p *NoopProtocol) GetShareableKey() string {
	return ""
}

// Set the encryption key
func (p *NoopProtocol) SetEncryptionKey(key string) {
}

// Get the version of the protocol
func (p *NoopProtocol) Version() string {
	return "0.1.0"
}

// Create a new handshake protocol
func NewNoOpProtocol() (*NoopProtocol, error) {
	return &NoopProtocol{}, nil
}
