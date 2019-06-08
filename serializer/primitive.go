package serializer

import "io"

// Primitive ...
type Primitive interface {
	Read(r io.Reader) ([]byte, error)
	Write(i []byte, w io.Writer) error
}

// EncodeInt ...
func EncodeInt(length int, encoded uint64) error {
	// To avoid reallocations, grow capacity to the largest possible size for this integer
	bb := make([]byte, 0, length)

	if encoded == 0 {
		bb = append(bb, byte(0))
		return nil
	}

	for encoded > 0 {
		b := byte(encoded & 127)
		encoded = encoded >> 7
		if !(encoded == 0) {
			b |= 128
		}

		bb = append(bb, b)
	}

	return nil
}

// WriteLong ...
// 	- TODO: include reference to spec
// 	- TMP: Writes the expected length of the upcomming message
// https://avro.apache.org/docs/1.8.1/spec.html
func WriteLong(r int64, w io.Writer) error {
	const maxByteSize = 10

	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))

	err := EncodeInt(maxByteSize, encoded)
	if err != nil {
		return err
	}

	return nil
}
