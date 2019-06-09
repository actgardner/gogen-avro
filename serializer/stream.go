package serializer

import "io"

// Stream low level Reader, Writer implementation of the ongoing data stream.
type Stream struct {
	Writer io.Writer
	Reader io.Reader
}
