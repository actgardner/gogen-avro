package serializer

import (
	"io"
)

// EncodeInt encodes the given interger using variable-length zig-zag coding.
// https://avro.apache.org/docs/1.8.1/spec.html#binary_encoding
func EncodeInt(length int, i uint64) []byte {
	// To avoid reallocations, grow capacity to the largest possible size for this integer
	bb := make([]byte, 0, length)

	if i == 0 {
		bb = append(bb, byte(0))
		return bb
	}

	for a := 0; i > 0; a++ {
		b := byte(i & 127)
		i = i >> 7
		if !(i == 0) {
			b |= 128
		}

		bb = append(bb, b)
	}

	return bb
}

// WriteInt writes the given int to the underlaying data stream.
func WriteInt(w io.Writer, r int32) error {
	const maxByteSize = 5

	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r>>downShift))

	bb := EncodeInt(maxByteSize, encoded)
	_, err := w.Write(bb)

	return err
}

// ReadInt interperates the next byte of the underlaying data stream as a int.
func ReadInt(r io.Reader) (int32, error) {
	var v uint32
	buf := make([]byte, 1)

	for shift := uint(0); ; shift += 7 {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return 0, err
		}

		b := buf[0]
		v |= uint32(b&127) << shift

		if b&128 == 0 {
			break
		}
	}

	i := (int32(v>>1) ^ -int32(v&1))
	return i, nil
}

// WriteMapInt interperates the next bytes of the underlaying data stream as a map[string]int
func WriteMapInt(w io.Writer, m map[string]int32) error {
	err := WriteMessageLength(w, int64(len(m)))
	if err != nil || len(m) == 0 {
		return err
	}

	for key, val := range m {
		err = WriteString(w, key)
		if err != nil {
			return err
		}

		err = WriteInt(w, val)
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

// ReadMapInt interperates the next bytes of the underlaying data stream as a map[string]int
func ReadMapInt(r io.Reader) (map[string]int32, error) {
	m := make(map[string]int32)

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

			val, err := ReadInt(r)
			if err != nil {
				return nil, err
			}

			m[key] = val
		}
	}

	return m, nil
}
