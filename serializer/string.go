package serializer

import "io"

// WriteString writes the given string to the given data stream.
func WriteString(w io.Writer, i string) error {
	return WriteByte(w, []byte(i))
}

// ReadString reads the next length header and message block of len(m) and interperates it as a string
func ReadString(r io.Reader) (string, error) {
	bb, err := ReadByte(r)
	if err != nil {
		return "", err
	}

	return string(bb), nil
}

// WriteMapString writes the given map[string]string to the given data stream.
func WriteMapString(w io.Writer, m map[string]string) error {
	err := WriteMessageLength(w, int64(len(m)))
	if err != nil || len(m) == 0 {
		return err
	}

	for key, val := range m {
		err = WriteString(w, key)
		if err != nil {
			return err
		}

		err = WriteString(w, val)
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

// ReadMapString interperates the next bytes of the given data stream as a map[string]string
func ReadMapString(r io.Reader) (map[string]string, error) {
	m := make(map[string]string)

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

			val, err := ReadString(r)
			if err != nil {
				return nil, err
			}

			m[key] = val
		}
	}

	return m, nil
}
