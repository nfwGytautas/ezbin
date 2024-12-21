package requests

type HandshakeRequest struct {
	UserIdentifier string `json:"userIdentifier"`
}

type HandshakeResponse struct {
	Okay      bool   `json:"okay"`
	Framesize int    `json:"framesize"`
	Protocol  string `json:"protocol"`
}
