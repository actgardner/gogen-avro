package serializer

import "io"

// EncodeFloat encodes the given float using variable-length zig-zag coding.
// https://avro.apache.org/docs/1.8.1/spec.html#binary_encoding
func EncodeFloat(w io.Writer, length int, bits uint64) error {
	bb := make([]byte, length)

	for i := 0; i < length; i++ {
		bb[i] = byte(bits & 255)
		bits = bits >> 8
	}

	_, err := w.Write(bb)
	return err
}
