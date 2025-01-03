package requests

type HandshakeRequest struct {
	UserIdentifier string `json:"id"`
	Key            string `json:"key"`
}

type HandshakeResponse struct {
	Okay      bool `json:"okay"`
	Framesize int  `json:"framesize"`
}
