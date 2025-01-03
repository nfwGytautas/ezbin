package protocol

// RSA protocol implementation
type RSAProtocol struct {
}

// Get the name of the protocol
func (p *RSAProtocol) Name() string {
	return PROTOCOL_RSA
}

// Get the version of the protocol
func (p *RSAProtocol) Version() string {
	return "0.1.0"
}
