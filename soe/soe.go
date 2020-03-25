package soe

import (
	"encoding/binary"
	"errors"
	"io"
)

var HeaderV1 = []byte{0xC3, 0x01}

func WriteHeader(w io.Writer, header []byte) error {
	fp := HeaderV1
	err := binary.Write(w, binary.LittleEndian, fp)
	if err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, header)
}

func ReadHeader(r io.Reader) ([]byte, error) {
	b := make([]byte, 2)
	n, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if b[0] != HeaderV1[0] || b[1] != HeaderV1[1] {
		return nil, errors.New("avro v1 header invalid")
	}

	var header []byte

	header = make([]byte, 8)

	n, err = r.Read(header)

	if n != 8 || err != nil {
		return nil, errors.New("fingerprint header not read")
	}

	return header, nil
}
