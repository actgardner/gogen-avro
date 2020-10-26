package main

import (
	"bytes"
	"io"
)

type Field interface {
	SetBoolean(v bool)
}

type bField bool

func (b *bField) SetBoolean(v bool) {
	*b = bField(v)
}

func assignBool(r io.Reader, f Field) error {
	v, err := readBool(r)
	if err != nil {
		return err
	}
	f.SetBoolean(v)
	return nil
}

func readBool(r io.Reader) (bool, error) {
	var b byte
	var err error
	bs := make([]byte, 1)
	_, err = io.ReadFull(r, bs)
	if err != nil {
		return false, err
	}
	b = bs[0]
	return b == 1, nil
}

func callAssign(r io.Reader, f Field) error {
	return assignBool(r, f)
}

func main() {
	buf := bytes.NewBuffer([]byte{0x01})
	var b bField
	callAssign(buf, &b)
}
