package connection

import (
	"encoding/json"
	"io"
	"net"
	"strings"

	"github.com/nfwGytautas/ezbin/ezbin"
)

// Frame is a struct that represents a ezbin tcp frame
type Frame struct {
	conn   net.Conn
	buffer []byte

	lastReadBytes int
	writeSize     int
}

// NewFrame creates a new frame
func NewFrame(conn net.Conn, buffer []byte) *Frame {
	return &Frame{
		conn:   conn,
		buffer: buffer,
	}
}

// Read the frame from the connection
func (f *Frame) Read() error {
	n, err := f.conn.Read(f.buffer)
	if err != nil {
		return err
	}

	f.lastReadBytes = n
	f.writeSize = 0

	return nil
}

// Write the frame to the connection
func (f *Frame) Write() error {
	if f.writeSize == 0 {
		return ezbin.ErrNothingToWrite
	}

	_, err := f.conn.Write(f.buffer[:f.writeSize])
	if err != nil {
		return err
	}

	f.lastReadBytes = 0
	f.writeSize = 0

	return nil
}

// FromJSON encodes a JSON object into a frame
func (f *Frame) FromJSON(in interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	err = f.writeRange(HEADER_SIZE_BYTES, data)
	if err != nil {
		return err
	}

	return nil
}

// ToJSON decodes a frame into a JSON object
func (f *Frame) ToJSON(out interface{}) error {
	err := json.Unmarshal(f.buffer[HEADER_SIZE_BYTES:f.lastReadBytes], out)
	if err != nil {
		return err
	}

	return nil
}

// Get the header of the frame
func (f *Frame) GetHeader() string {
	header := f.buffer[:HEADER_SIZE_BYTES]
	return strings.TrimRight(string(header), "\x00")
}

// Write the header to the frame
func (f *Frame) SetHeader(header string) error {
	if len(header) > HEADER_SIZE_BYTES {
		return ezbin.ErrHeaderTooLarge
	}

	headerBin := strings.Repeat("\x00", HEADER_SIZE_BYTES)
	headerBin = header + headerBin[len(header):]

	// Write the header
	err := f.writeRange(0, []byte(headerBin))
	if err != nil {
		return err
	}

	return nil
}

// Write to a frame from a generic writer interface
func (f *Frame) TransferToWriter(w io.Writer) error {
	_, err := w.Write(f.buffer[HEADER_SIZE_BYTES:f.lastReadBytes])
	if err != nil {
		return err
	}

	return nil
}

// Read from a frame to a generic reader interface
func (f *Frame) TransferFromReader(r io.Reader, start int) error {
	if start < 0 {
		return ezbin.ErrInvalidStart
	}

	n, err := r.Read(f.buffer[HEADER_SIZE_BYTES+start:])
	if err != nil {
		return err
	}

	if n+start > f.writeSize {
		f.writeSize = HEADER_SIZE_BYTES + n + start
	}

	return nil
}

// Get the size of the frame without header
func (f *Frame) GetFrameSize() int {
	return f.writeSize - HEADER_SIZE_BYTES
}

// Get the number of bytes read
func (f *Frame) GetNumReadBytes() int {
	return f.lastReadBytes - HEADER_SIZE_BYTES
}

// Write a range to the frame
func (f *Frame) writeRange(start int, data []byte) error {
	if len(data) > len(f.buffer)-start {
		return ezbin.ErrBufferTooSmall
	}

	if len(data) == 0 {
		return nil
	}

	copy(f.buffer[start:], data)

	if start+len(data) > f.writeSize {
		f.writeSize = start + len(data)
	}

	return nil
}
