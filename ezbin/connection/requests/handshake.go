package requests

type HandshakeRequest struct {
	UserIdentifier string `json:"userIdentifier"`
	Protocol       string `json:"protocol"`
	Key            string `json:"key"`
}

type HandshakeResponse struct {
	Okay      bool   `json:"okay"`
	Framesize int    `json:"framesize"`
	Key       string `json:"key"`
}
