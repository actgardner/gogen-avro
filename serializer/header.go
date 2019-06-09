package serializer

import "io"

// HeaderWriter writes header messages to the underlaying data stream.
type HeaderWriter struct {
	writer io.Writer
}

// WriteMessageLength writes ammount of bytes that could be expected by a consumer for the upcomming message.
// https://avro.apache.org/docs/1.8.1/spec.html#Message+Framing
func (h HeaderWriter) WriteMessageLength(r int64) error {
	const maxByteSize = 10

	downShift := uint64(63)
	length := uint64((r << 1) ^ (r >> downShift))

	encoded := EncodeInt(maxByteSize, length)
	_, err := h.writer.Write(encoded)

	return err
}
