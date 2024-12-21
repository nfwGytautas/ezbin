package shared

import "errors"

// Write a binary sub range
func WriteSubRange(buffer []byte, start int, data []byte) error {
	if len(data) > len(buffer)-start {
		return errors.New("data exceeds buffer size")
	}

	if len(data) == 0 {
		return nil
	}

	copy(buffer[start:], data)

	return nil
}
