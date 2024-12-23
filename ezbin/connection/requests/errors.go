package requests

import "errors"

var (
	ErrExceededFrameSize = errors.New("exceeded frame size")
)
