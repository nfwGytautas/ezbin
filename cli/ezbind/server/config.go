package server

import (
	"github.com/nfwGytautas/ezbin/ezbin/protocol"
	ezbin_server "github.com/nfwGytautas/ezbin/ezbin/server"
)

// Create a new default peer config
func NewPeerConfig() (*ezbin_server.DaemonConfig, error) {
	dc, err := ezbin_server.NewDefaultDaemonConfig()
	if err != nil {
		return nil, err
	}

	// Add peer config
	protocols, err := protocol.GenerateProtocols()
	if err != nil {
		return nil, err
	}

	dc.Peer = &struct {
		Protocol protocol.Protocols "yaml:\"protocol\""
	}{
		Protocol: protocols,
	}

	return dc, nil
}
