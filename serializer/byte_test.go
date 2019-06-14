package serializer

import (
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestReadingBytes(t *testing.T) {
	inputs := map[string][]byte{
		"never": []byte{10, 110, 101, 118, 101, 114},
		"gonna": []byte{10, 103, 111, 110, 110, 97},
		"give":  []byte{8, 103, 105, 118, 101},
		"you":   []byte{6, 121, 111, 117},
		"up":    []byte{4, 117, 112},
	}

	r, w := io.Pipe()

	for expected, input := range inputs {
		go w.Write(input)

		result, err := ReadByte(r)
		if err != nil {
			t.Fatal(err)
		}

		if string(result) != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %s recieved %s", input, expected, result)
		}
	}
}

func BenchmarkWritingBytes(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	inputs := make([][]byte, b.N)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, []byte(RandStringRunes(100)))
	}

	r, w := io.Pipe()
	go ioutil.ReadAll(r)

	b.ResetTimer()

	for _, input := range inputs {
		WriteByte(w, input)
	}
}
