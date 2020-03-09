package canonical

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvroCRC64Fingerprint(t *testing.T) {
	cases := []struct {
		schema      string
		fingerprint string
	}{
		{`"int"`, "8f5c393f1ad57572"},
		{`"float"`, "90d7a83ecb027c4d"},
		{`"long"`, "b71df49344e154d0"},
		{`"double"`, "7e95ab32c035758e"},
		{`"bytes"`, "651920c3da16c04f"},
	}

	for _, c := range cases {
		assert.Equal(t, c.fingerprint, AvroCRC64Fingerprint([]byte(c.schema)))
	}
}
