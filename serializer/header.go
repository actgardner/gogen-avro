package serializer

import "io"

// NewHeader constructs a new header processer for the given stream
func NewHeader(stream Stream) Header {
	h := Header{
		Stream: stream,
	}

	return h
}

// Header writes header messages to the underlaying data stream.
type Header struct {
	Stream
}

// WriteMessageLength writes ammount of bytes that could be expected by a consumer for the upcomming message.
// Any error encountered while writing the message header is returned.
// https://avro.apache.org/docs/1.8.1/spec.html#Message+Framing
func (h Header) WriteMessageLength(r int64) error {
	const maxByteSize = 10

	downShift := uint64(63)
	length := uint64((r << 1) ^ (r >> downShift))

	encoded := EncodeInt(maxByteSize, length)
	_, err := h.Stream.Write(encoded)

	return err
}

// ReadMessageLength reads the next message header which contains the expecting message length.
// It returns the total byte length of the upcomming message and any error encountered.
// https://avro.apache.org/docs/1.8.1/spec.html#Message+Framing
func (h Header) ReadMessageLength() (int64, error) {
	var l uint64
	buf := make([]byte, 1)

	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(h.Stream, buf); err != nil {
			return 0, err
		}

		b := buf[0]
		l |= uint64(b&127) << shift
		if b&128 == 0 {
			break
		}
	}

	length := (int64(l>>1) ^ -int64(l&1))
	return length, nil
}
