package vm

import (
	"bytes"
	"testing"
)

func BenchmarkReadBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer([]byte{1})
		readBool(r)
	}
}
