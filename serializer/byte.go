package serializer

import (
	"io"
)

// NewByte constructs a new byte processer for the given stream
func NewByte(stream Stream) Byte {
	b := Byte{
		Header: NewHeader(stream),
		Stream: stream,
	}

	return b
}

// Byte low level byte reader, writer implementation.
// And Read, Writer implementation of the byte primitive.
type Byte struct {
	Header
	Stream
}

// Write writes the given byte buffer and the expecting message length to the underlaying data stream.
func (s *Byte) Write(i []byte) error {
	err := s.WriteMessageLength(int64(len(i)))
	if err != nil {
		return err
	}

	_, err = s.Stream.Write(i)
	return err
}

// ReadNext reads the next length header and message block of len(m)
func (s *Byte) ReadNext() ([]byte, error) {
	length, err := s.ReadMessageLength()
	if err != nil {
		return nil, err
	}

	bb := make([]byte, length)
	_, err = io.ReadFull(s.Stream, bb)
	if err != nil {
		return nil, err
	}

	return bb, nil
}
