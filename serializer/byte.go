package serializer

import (
	"io"
)

// WriteByte writes the given byte buffer and the expecting message length to the underlaying data stream.
func WriteByte(w io.Writer, i []byte) error {
	err := WriteMessageLength(w, int64(len(i)))
	if err != nil {
		return err
	}

	_, err = w.Write(i)
	return err
}

// ReadByte reads the next byte message block of len(m)
func ReadByte(r io.Reader) ([]byte, error) {
	length, err := ReadMessageLength(r)
	if err != nil {
		return nil, err
	}

	bb := make([]byte, length)
	_, err = r.Read(bb)
	if err != nil {
		return nil, err
	}

	return bb, nil
}
