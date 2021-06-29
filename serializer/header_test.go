package serializer

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestReadingHeaderMessageLength(t *testing.T) {
	inputs := map[int64][]byte{
		4:  []byte{8},
		1:  []byte{2},
		10: []byte{20},
	}

	for expected, input := range inputs {
		bb := bytes.NewBuffer(input)

		result, err := ReadMessageLength(bb)
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %d recieved %d", input, expected, result)
		}

		if bb.Len() != 0 {
			t.Fatal("not all bytes have been read from the byte buffer")
		}
	}
}

func TestWritingHeaderMessageLength(t *testing.T) {
	inputs := map[int64][]byte{
		4:  []byte{8},
		1:  []byte{2},
		10: []byte{20},
	}

	for input, expected := range inputs {
		r, w := io.Pipe()

		go func() {
			err := WriteMessageLength(w, input)
			if err != nil {
				t.Fatal(err)
			}

			w.Close()
		}()

		bb, _ := ioutil.ReadAll(r)

		if len(bb) != len(expected) {
			t.Fatalf("the returned byte buffer has an unexpected length: %b, %b\n", bb, expected)
		}

		for i, b := range bb {
			if b != expected[i] {
				t.Fatalf("unexpected byte encountered: %v, %v\n", b, expected[i])
			}
		}
	}
}

func BenchmarkReadingMessageLength(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	input := make([]byte, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		input = append(input, EncodeInt(10, uint64(rand.Intn(10000)))...)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadMessageLength(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingMessageLength(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	inputs := make([]int64, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, int64(rand.Intn(10000)))
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteMessageLength(bb, input)
	}
}
