package soe

import (
	"bytes"
	"encoding/binary"
	"github.com/actgardner/gogen-avro/schema/canonical"
	"io"
)

type soWriter struct {
	w      *bytes.Buffer
	header []byte
}

func NewWriter(buf *bytes.Buffer, header []byte) *soWriter {
	err := avroVersionHeader(buf, header)
	if err != nil {
		return nil
	}
	return &soWriter{w: buf, header: header}
}

func (s *soWriter) Write(p []byte) (n int, err error) {
	return s.w.Write(p)
}

func (s *soWriter) Grow(n int) {
	s.w.Grow(n)
}

func (s *soWriter) WriteByte(c byte) error {
	return s.w.WriteByte(c)
}

func avroVersionHeader(writer io.Writer, header []byte) error {
	fp := canonical.HeaderV1
	err := binary.Write(writer, binary.LittleEndian, fp)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, header)
	return err
}
