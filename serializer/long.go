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
