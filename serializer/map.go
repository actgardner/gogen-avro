package serializer

import (
	"io"
)

// ReadLengthNextMapBlock reads the length of the next map block
func ReadLengthNextMapBlock(r io.Reader) (block int64, err error) {
	block, err = ReadMessageLength(r)
	if err != nil {
		return block, err
	}

	if block == 0 {
		return block, io.EOF
	}

	if block < 0 {
		block = -block
		_, err := ReadMessageLength(r)
		if err != nil {
			return block, err
		}
	}

	return block, nil
}
