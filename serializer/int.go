package serializer

import "io"

// EncodeInt encodes the given interger using variable-length zig-zag coding.
// https://avro.apache.org/docs/1.8.1/spec.html#binary_encoding
func EncodeInt(length int, i uint64) []byte {
	// To avoid reallocations, grow capacity to the largest possible size for this integer
	bb := make([]byte, 0, length)

	if i == 0 {
		bb = append(bb, byte(0))
		return bb
	}

	for i > 0 {
		b := byte(i & 127)
		i = i >> 7
		if !(i == 0) {
			b |= 128
		}

		bb = append(bb, b)
	}

	return bb
}

// NewInt constructs a new int processer for the given stream
func NewInt(stream Stream) Int {
	i := Int{
		Stream: stream,
	}

	return i
}

// Int Read, Write implementation of the int (int32) primitive.
type Int struct {
	Stream
}

// Read interperates the next byte of the underlaying data stream as a int.
func (i *Int) Read() (int32, error) {
	var v uint32
	buf := make([]byte, 1)

	for shift := uint(0); ; shift += 7 {
		_, err := io.ReadFull(i.Stream, buf)
		if err != nil {
			return 0, err
		}

		b := buf[0]
		v |= uint32(b&127) << shift

		if b&128 == 0 {
			break
		}
	}

	r := (int32(v>>1) ^ -int32(v&1))
	return r, nil
}

// Write writes the given int to the underlaying data stream.
func (i *Int) Write(r int32) error {
	const maxByteSize = 5

	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r>>downShift))

	bb := EncodeInt(maxByteSize, encoded)
	_, err := i.Stream.Write(bb)

	return err
}
