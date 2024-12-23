package protocol

// `ecdsa` protocol implementation
type ECDSAProtocol struct {
}

// Get the name of the protocol
func (p *ECDSAProtocol) Name() string {
	return "ECDSA"
}

// Get the version of the protocol
func (p *ECDSAProtocol) Version() string {
	return "0.1.0"
}

// Generate new data for the protocol
func (p *ECDSAProtocol) GenerateNew() (ProtocolData, error) {
	return ProtocolData{
		"version":    "0.1.0",
		"publicKey":  "ecdsa-public-key",
		"privateKey": "ecdsa-private-key",
	}, nil
}
