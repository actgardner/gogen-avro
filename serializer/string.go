package serializer

// NewString constructs a new string processer for the given stream
func NewString(stream Stream) String {
	s := String{
		Byte: NewByte(stream),
	}

	return s
}

// String Read, Writer implementation of the string primitive.
type String struct {
	Byte
}

// Write writes the given byte buffer and the expecting message length to the underlaying data stream.
func (s *String) Write(i []byte) error {
	return s.Byte.Write(i)
}

// Read reads the next length header and message block of len(m)
func (s *String) Read() (string, error) {
	bb, err := s.Byte.Read()
	if err != nil {
		return "", err
	}

	return string(bb), nil
}
