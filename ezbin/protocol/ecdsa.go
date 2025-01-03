package protocol

// `ecdsa` protocol implementation
type ECDSAProtocol struct {
}

// Get the name of the protocol
func (p *ECDSAProtocol) Name() string {
	return PROTOCOL_ECDSA
}

// Get the version of the protocol
func (p *ECDSAProtocol) Version() string {
	return "0.1.0"
}
