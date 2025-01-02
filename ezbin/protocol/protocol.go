package protocol

// Interface for `ezbin` data sending and receiving protocols
type SRProtocol interface {
	// Get the name of the protocol
	Name() string

	// Get the version of the protocol
	Version() string

	// Generate new data for the protocol
	GenerateNew() (ProtocolData, error)
}

type ProtocolData = map[string]string

// Get all supported protocols
func GetSupportedProtocols() []SRProtocol {
	return []SRProtocol{
		&ECDSAProtocol{},
		&RSAProtocol{},
		&NoopProtocol{},
	}
}

// Get a protocol by name
func GetProtocolByName(name string) SRProtocol {
	for _, p := range GetSupportedProtocols() {
		if p.Name() == name {
			return p
		}
	}
	return nil
}
