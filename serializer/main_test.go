package serializer

import "io"

// NewStream sets up a io.Pipe for streaming tests
func NewStream() Stream {
	r, w := io.Pipe()
	s := Stream{
		Reader: r,
		Writer: w,
	}

	return s
}
