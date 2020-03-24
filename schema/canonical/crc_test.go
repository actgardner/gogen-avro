package canonical

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAvroCRC64Fingerprint(t *testing.T) {
	cases := []struct {
		schema      string
		fingerprint string
	}{
		{`"long"`, "b71df49344e154d0"},
		{`"double"`, "7e95ab32c035758e"},
		{`"bytes"`, "651920c3da16c04f"},
	}

	for _, c := range cases {
		output := AvroCRC64Fingerprint([]byte(c.schema))
		assert.Equal(t, c.fingerprint, hex.EncodeToString(output))
	}
}
