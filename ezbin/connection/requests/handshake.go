package requests

type HandshakeRequest struct {
	// User identifier
	UserIdentifier string `json:"userIdentifier"`
}

type HandshakeResponse struct {
}

// Get the request header
func (r HandshakeRequest) Header() string {
	return "HANDSHAKE"
}
