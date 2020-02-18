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
		{`{"type": "int"}`, "8f5c393f1ad57572"},
	}

	for _, c := range cases {
		assert.Equal(t, c.fingerprint, AvroCRC64Fingerprint([]byte(c.schema)))
	}
}
