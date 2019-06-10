package serializer

import (
	"io"
)

// NewLong constructs a new long int processer for the given stream
func NewLong(stream Stream) Long {
	l := Long{
		Stream: stream,
	}

	return l
}

// Long Read, Write implementation of the long int (int64) primitive
type Long struct {
	Stream
}

// Read interperates the next byte of the underlaying data stream as a long int.
func (l *Long) Read() (int64, error) {
	var v uint64
	buf := make([]byte, 1)

	for shift := uint(0); ; shift += 7 {
		_, err := io.ReadFull(l.Stream, buf)
		if err != nil {
			return 0, err
		}

		b := buf[0]
		v |= uint64(b&127) << shift

		if b&128 == 0 {
			break
		}
	}

	r := (int64(v>>1) ^ -int64(v&1))
	return r, nil
}

// Write writes the given long int to the underlaying data stream.
func (l *Long) Write(i int64) error {
	const maxByteSize = 10

	downShift := uint64(63)
	encoded := uint64((i << 1) ^ (i >> downShift))

	bb := EncodeInt(maxByteSize, encoded)
	_, err := l.Stream.Write(bb)

	return err
}
