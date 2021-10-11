package canonical

import (
	"testing"
)

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
)

func TestAvroCRC64Fingerprint(t *testing.T) {
	cases := []struct {
		schema      string
		fingerprint string
	}{
		{`"int"`, "8f5c393f1ad57572"},
		{`"long"`, "b71df49344e154d0"},
		{`"double"`, "7e95ab32c035758e"},
		{`"bytes"`, "651920c3da16c04f"},
	}

	for _, c := range cases {
		output := AvroCRC64Fingerprint([]byte(c.schema))
		assert.Equal(t, c.fingerprint, hex.EncodeToString(output))
	}
}
