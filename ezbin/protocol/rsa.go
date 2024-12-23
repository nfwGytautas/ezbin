package protocol

// RSA protocol implementation
type RSAProtocol struct {
}

// Get the name of the protocol
func (p *RSAProtocol) Name() string {
	return "RSA"
}

// Get the version of the protocol
func (p *RSAProtocol) Version() string {
	return "0.1.0"
}

// Generate new data for the protocol
func (p *RSAProtocol) GenerateNew() (ProtocolData, error) {
	return ProtocolData{
		"version":    "0.1.0",
		"publicKey":  "rsa-public-key",
		"privateKey": "rsa-private-key",
	}, nil
}
