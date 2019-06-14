package serializer

import (
	"io"
	"io/ioutil"
	"math/rand"
	"time"

	"testing"
)

func TestReadingString(t *testing.T) {
	inputs := map[string][]byte{
		"john": []byte{8, 106, 111, 104, 110},
		"doe":  []byte{6, 100, 111, 101},
	}

	r, w := io.Pipe()

	for expected, input := range inputs {
		go w.Write(input)

		result, err := ReadString(r)
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %s recieved %s", input, expected, result)
		}
	}
}

func TestWritingString(t *testing.T) {
	inputs := map[string][]byte{
		"john": []byte{8, 106, 111, 104, 110},
		"doe":  []byte{6, 100, 111, 101},
	}

	for input, expected := range inputs {
		r, w := io.Pipe()

		go func() {
			err := WriteString(w, input)
			if err != nil {
				t.Fatal(err)
			}

			w.Close()
		}()

		bb, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}

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

func BenchmarkWritingString(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	inputs := make([]string, b.N)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, RandStringRunes(100))
	}

	r, w := io.Pipe()
	go ioutil.ReadAll(r)

	b.ResetTimer()

	for _, input := range inputs {
		WriteString(w, input)
	}
}

func BenchmarkReadingString(b *testing.B) {
	r, w := io.Pipe()

	go func() {
		for i := 0; i < b.N; i++ {
			value := RandStringRunes(100)
			WriteString(w, value)
		}

		w.Close()
	}()

	bb, _ := ioutil.ReadAll(r)
	r, w = io.Pipe()
	go w.Write(bb)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadString(r)
		if err != nil {
			b.Fatal(err)
		}
	}
}
