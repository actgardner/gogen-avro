package serializer

import (
	"io"
)

// StringWriter Writer implementation of the string primitive.
type StringWriter struct {
	HeaderWriter
	writer io.Writer
}

// Write writes the given byte buffer and the expecting message length to the underlaying data stream.
func (s *StringWriter) Write(i []byte) error {
	err := s.WriteMessageLength(int64(len(i)))
	if err != nil {
		return err
	}

	_, err = s.writer.Write(i)
	return err
}
