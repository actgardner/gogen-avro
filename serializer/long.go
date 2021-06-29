package serializer

import (
	"io"
)

// ReadLong interperates the next byte of the given data stream as a long int.
func ReadLong(r io.Reader) (int64, error) {
	var v uint64
	buf := make([]byte, 1)

	for shift := uint(0); ; shift += 7 {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return 0, err
		}

		b := buf[0]
		v |= uint64(b&127) << shift

		if b&128 == 0 {
			break
		}
	}

	l := (int64(v>>1) ^ -int64(v&1))
	return l, nil
}

// WriteLong writes the given long int to the given data stream.
func WriteLong(w io.Writer, i int64) error {
	const maxByteSize = 10

	downShift := uint64(63)
	encoded := uint64((i << 1) ^ (i >> downShift))

	bb := EncodeInt(maxByteSize, encoded)
	_, err := w.Write(bb)

	return err
}

// WriteMapLong writes the given map[string]int64 to the given data stream.
func WriteMapLong(w io.Writer, m map[string]int64) error {
	err := WriteMessageLength(w, int64(len(m)))
	if err != nil || len(m) == 0 {
		return err
	}

	for key, val := range m {
		err = WriteString(w, key)
		if err != nil {
			return err
		}

		err = WriteLong(w, val)
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

// ReadMapLong interperates the next bytes of the given data stream as a map[string]int64
func ReadMapLong(r io.Reader) (map[string]int64, error) {
	m := make(map[string]int64)

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

			val, err := ReadLong(r)
			if err != nil {
				return nil, err
			}

			m[key] = val
		}
	}

	return m, nil
}
