package soe

import (
	"bytes"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAvroVersionHeader(t *testing.T) {
	cases := []struct {
		header   []byte
		expected string
	}{
		{[]byte{0x0c, 0x94, 0x7f, 0x60, 0x1d, 0xe7, 0xce, 0x84}, "c3010c947f601de7ce84"},
	}

	for _, c := range cases {
		b := make([]byte, 0, 2)
		output := bytes.NewBuffer(b)
		err := avroVersionHeader(output, c.header)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, hex.EncodeToString(output.Bytes()))
	}
}
