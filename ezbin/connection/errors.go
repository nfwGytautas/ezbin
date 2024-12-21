package connection

import "errors"

var (
	ErrConnectionFailed = errors.New("connection to peer failed")
	ErrHeaderTooLarge   = errors.New("header too large")
	ErrUnknownHeader    = errors.New("unknown header")
)
