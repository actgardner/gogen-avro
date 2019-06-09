package serializer

import "io"

// Reader is the interface that wraps the basic Read method.
//
// ReadNext reads the next block of bytes from the underlaying data stream.
// It returns the bytes read inside a buffer and any error encountered.
//
// Implementations define the type of the returned byte buffer.
type Reader interface {
	ReadNext() ([]byte, error)
}

// Writer is the interface that wraps the basic Write method.
//
// Write writes len(p) bytes from p to the underlying data stream.
// A expected length header of len(p) is written before the actual message.
// It returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p). Write must not modify the slice data, even temporarily.
//
// Implementations define the type of p.
type Writer interface {
	Write(p []byte) error
}

// Closer The behavior of Close after the first call is undefined. Specific implementations may document their own behavior.
type Closer interface {
	Close() error
}

// Primitive is the interface that groups the basic ReadNext, Write and Close methods.
type Primitive interface {
	Reader
	Writer
	Closer
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
