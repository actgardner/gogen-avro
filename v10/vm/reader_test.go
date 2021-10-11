package vm

import (
	"bytes"
	"testing"
)

func BenchmarkReadBool(b *testing.B) {
	buf := make([]byte, 8)
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer([]byte{1})
		readBool(r, buf)
	}
}
