package ezbin

import "github.com/nfwGytautas/ezbin/ezbin/protocol"

// Get all supported protocols
func GetSupportedProtocols() []protocol.SRProtocol {
	return []protocol.SRProtocol{
		&protocol.ECDSAProtocol{},
		&protocol.RSAProtocol{},
		&protocol.NoopProtocol{},
	}
}

// Get a protocol by name
func GetProtocolByName(name string) protocol.SRProtocol {
	for _, p := range GetSupportedProtocols() {
		if p.Name() == name {
			return p
		}
	}
	return nil
}
