package serializer

import (
	"bytes"
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

func TestWritingBytes(t *testing.T) {
	inputs := map[string][]byte{
		"john": []byte{8, 106, 111, 104, 110},
		"doe":  []byte{6, 100, 111, 101},
	}

	for input, expected := range inputs {
		r, w := io.Pipe()

		go func() {
			err := WriteByte(w, []byte(input))
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
				t.Fatalf("unexpected byte encountered: %b, %b\n", b, expected[i])
			}
		}
	}
}

func BenchmarkReadingBytes(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		value := RandStringRunes(100)
		WriteString(bb, value)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadByte(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingBytes(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	inputs := make([][]byte, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, []byte(RandStringRunes(100)))
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteByte(bb, input)
	}
}
