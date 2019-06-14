package serializer

import (
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestHeaderMessageLength(t *testing.T) {
	inputs := map[int64][]byte{
		4:  []byte{8},
		1:  []byte{2},
		10: []byte{20},
	}

	r, w := io.Pipe()

	for expected, input := range inputs {
		go w.Write(input)

		result, err := ReadMessageLength(r)
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %d recieved %d", input, expected, result)
		}
	}
}

func BenchmarkWritingMessageLength(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	inputs := make([]int64, b.N)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, int64(rand.Intn(10000)))
	}

	r, w := io.Pipe()
	go ioutil.ReadAll(r)

	b.ResetTimer()

	for _, input := range inputs {
		WriteMessageLength(w, input)
	}
}
