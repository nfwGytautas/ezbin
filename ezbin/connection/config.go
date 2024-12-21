package connection

import "time"

// The size of the header in bytes
const HEADER_SIZE_BYTES = 16

// The amount of bytes that is guaranteed to exist in the tcp packet from the start,
// can grow depending on the peer settings but cannot be less
const HANDSHAKE_BUFFER_SIZE = 1024

// Timeout in seconds for handshake request
const HANDSHAKE_READ_TIMEOUT = time.Second * 5
