package ezbin

import "errors"

var (
	ErrConnectionFailed     = errors.New("connection to peer failed")
	ErrHeaderTooLarge       = errors.New("header too large")
	ErrUnknownHeader        = errors.New("unknown header")
	ErrIncorrectHeader      = errors.New("incorrect header")
	ErrHandshakeFailed      = errors.New("handshake failed")
	ErrHandshakeNotFinished = errors.New("handshake not finished")
	ErrBufferTooSmall       = errors.New("buffer too small")
	ErrUploadFailed         = errors.New("upload failed")
	ErrNothingToWrite       = errors.New("nothing to write")
	ErrInvalidStart         = errors.New("invalid start")
	ErrInvalidConfig        = errors.New("invalid config")
	ErrIdentityNotFound     = errors.New("identity not found")
	ErrPeerNotFound         = errors.New("peer not found")
	ErrPeerExists           = errors.New("peer already exists")
	ErrPackageNotFound      = errors.New("package not found")
	ErrExceededFrameSize    = errors.New("exceeded frame size")
	ErrUnknownProtocol      = errors.New("unknown protocol")
)
