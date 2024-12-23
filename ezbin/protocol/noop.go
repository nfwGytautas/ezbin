package protocol

// Noop protocol implementation
type NoopProtocol struct {
}

// Get the name of the protocol
func (p *NoopProtocol) Name() string {
	return "no-op"
}

// Get the version of the protocol
func (p *NoopProtocol) Version() string {
	return "0.1.0"
}

// Generate new data for the protocol
func (p *NoopProtocol) GenerateNew() (ProtocolData, error) {
	return ProtocolData{
		"version": "0.1.0",
	}, nil
}
