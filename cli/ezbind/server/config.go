package server

import (
	"github.com/nfwGytautas/ezbin/ezbin"
)

// Create a new default peer config
func NewPeerConfig() (*ezbin.DaemonConfig, error) {
	dc, err := ezbin.NewDefaultDaemonConfig()
	if err != nil {
		return nil, err
	}

	// Add peer config
	dc.Peer = &struct {
		Protocol string "yaml:\"protocol\""
	}{}

	return dc, nil
}
