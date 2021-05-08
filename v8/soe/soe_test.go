package soe

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRecord struct {
	fingerprint []byte
}

func (m *MockRecord) AvroCRC64Fingerprint() []byte {
	return m.fingerprint
}

func (m *MockRecord) Serialize(_ io.Writer) error {
	return nil
}

func TestAvroVersionHeader(t *testing.T) {
	cases := []struct {
		record   *MockRecord
		expected string
	}{
		{&MockRecord{[]byte{0x0c, 0x94, 0x7f, 0x60, 0x1d, 0xe7, 0xce, 0x84}}, "c3010c947f601de7ce84"},
	}

	for _, c := range cases {
		b := make([]byte, 0)
		output := bytes.NewBuffer(b)
		err := WriteRecord(output, c.record)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, hex.EncodeToString(output.Bytes()))
	}
}

func TestReadHeader(t *testing.T) {
	cases := []struct {
		header      []byte
		expected    []byte
		expectedLen int
	}{
		{[]byte{0xc3, 0x01, 0x7f, 0x60, 0x1d, 0xe7, 0xce, 0x84, 0x01, 0x02, 0x03}, []byte{0x7f, 0x60, 0x1d, 0xe7, 0xce, 0x84, 0x01, 0x02}, 1},
	}

	for _, c := range cases {
		input := bytes.NewBuffer(c.header)
		header, err := ReadHeader(input)
		assert.Nil(t, err)
		assert.Equal(t, c.expectedLen, len(input.Bytes()))
		assert.Equal(t, hex.EncodeToString(c.expected), hex.EncodeToString(header))
	}
}
