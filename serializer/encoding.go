package serializer

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
