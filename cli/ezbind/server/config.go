package server

import (
	ezbin_server "github.com/nfwGytautas/ezbin/ezbin/server"
)

// Create a new default peer config
func NewPeerConfig() (*ezbin_server.DaemonConfig, error) {
	dc, err := ezbin_server.NewDefaultDaemonConfig()
	if err != nil {
		return nil, err
	}

	// Add peer config
	dc.Peer = &struct {
		Protocol string "yaml:\"protocol\""
	}{}

	return dc, nil
}
