package singleobject

import (
	"bytes"
	"github.com/actgardner/gogen-avro/schema/canonical"
)

type soWriter struct {
	w      *bytes.Buffer
	header []byte
}

func NewWriter(buf *bytes.Buffer, header []byte) *soWriter {
	err := canonical.AvroVersionHeader(buf, header)
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
