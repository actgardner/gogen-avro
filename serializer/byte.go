package serializer

import (
	"io"
)

// WriteByte writes the given byte buffer and the expecting message length to the given data stream.
func WriteByte(w io.Writer, i []byte) error {
	err := WriteMessageLength(w, int64(len(i)))
	if err != nil {
		return err
	}

	_, err = w.Write(i)
	return err
}

// ReadByte reads the next byte message block of len(m) and interperates it as []byte
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

// WriteMapByte writes the given map[string][]byte to the given data stream.
func WriteMapByte(w io.Writer, m map[string][]byte) error {
	err := WriteMessageLength(w, int64(len(m)))
	if err != nil || len(m) == 0 {
		return err
	}

	for key, val := range m {
		err = WriteString(w, key)
		if err != nil {
			return err
		}

		err = WriteByte(w, val)
		if err != nil {
			return err
		}
	}

	// Mark the end of the map
	err = WriteMessageLength(w, 0)
	if err != nil {
		return err
	}

	return nil
}

// ReadMapByte interperates the next bytes of the given data stream as a map[string][]byte
func ReadMapByte(r io.Reader) (map[string][]byte, error) {
	m := make(map[string][]byte)

	for {
		block, err := ReadLengthNextMapBlock(r)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		for i := int64(0); i < block; i++ {
			key, err := ReadString(r)
			if err != nil {
				return nil, err
			}

			val, err := ReadByte(r)
			if err != nil {
				return nil, err
			}

			m[key] = val
		}
	}

	return m, nil
}
