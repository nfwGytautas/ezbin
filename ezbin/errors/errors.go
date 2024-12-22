package errors

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
)
