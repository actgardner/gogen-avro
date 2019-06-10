package serializer

import "io"

// Stream low level Reader, Writer implementation of a underlaing data stream.
type Stream struct {
	io.Reader
	io.Writer
}
