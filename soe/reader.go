package soe

import (
	"bytes"
	"errors"
)

type soReader struct {
	r      *bytes.Buffer
	header []byte
}

func NewReader(buf *bytes.Buffer) *soReader {
	header, err := readSOHeader(buf)
	if err != nil {
		return nil
	}
	return &soReader{r: buf, header: header}
}

func (s soReader) Read(p []byte) (int, error) {
	return s.r.Read(p)
}

func (s soReader) Bytes() []byte {
	return s.r.Bytes()
}

func readSOHeader(r *bytes.Buffer) ([]byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != 0xC3 {
		_ = r.UnreadByte()
		return nil, errors.New("header not read")
	}

	b, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	if b != 0x01 {
		_ = r.UnreadByte()
		return nil, errors.New("header not read")
	}

	var header []byte

	header = make([]byte, 8)

	n, err := r.Read(header)

	if n != 8 || err != nil {
		return nil, errors.New("header not read")
	}
	return header, nil
}
