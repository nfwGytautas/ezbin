package shared

import "errors"

func WriteSubRange(buffer []byte, start int, data []byte) error {
	if len(data) > len(buffer)-start {
		return errors.New("data exceeds buffer size")
	}

	if len(data) == 0 {
		return nil
	}

	for i := 0; i < len(data); i++ {
		buffer[start+i] = data[i]
	}

	return nil
}
