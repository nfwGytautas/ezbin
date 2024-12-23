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
