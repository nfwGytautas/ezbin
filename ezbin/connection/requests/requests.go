package requests

import (
	"strings"
)

// Arbitrary request to do by the server
type Request interface {
	// Get the request header
	Header() string
}

type Response interface {
}

// A command to send to the peer
type PeerCommand struct {
	Header  string  `json:"header"`
	Request Request `json:"request"`
}

// Convert a header to a request
func HeaderToRequestResponse(header string) (Request, Response) {
	switch strings.TrimRight(header, "\x00") {
	case HandshakeRequest{}.Header():
		return &HandshakeRequest{}, &HandshakeResponse{}
	}

	return nil, nil
}
