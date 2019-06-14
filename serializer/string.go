package serializer

import "io"

// WriteString writes the given byte buffer and the expecting message length to the underlaying data stream.
func WriteString(w io.Writer, i string) error {
	return WriteByte(w, []byte(i))
}

// ReadString reads the next length header and message block of len(m)
func ReadString(r io.Reader) (string, error) {
	bb, err := ReadByte(r)
	if err != nil {
		return "", err
	}

	return string(bb), nil
}
