package serializer

import (
	"bytes"
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

	for expected, input := range inputs {
		bb := bytes.NewBuffer(input)

		result, err := ReadString(bb)
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %s recieved %s", input, expected, result)
		}

		if bb.Len() != 0 {
			t.Fatal("not all bytes have been read from the byte buffer")
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
				t.Fatalf("unexpected byte encountered: %v, %v at index %d\n", b, expected[i], i)
			}
		}
	}
}

func BenchmarkWritingString(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	inputs := make([]string, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, RandStringRunes(100))
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteString(bb, input)
	}
}

func BenchmarkReadingString(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteString(bb, RandStringRunes(100))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadString(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestReadingMapString(t *testing.T) {
	type run struct {
		Input  []byte
		Output map[string]string
	}

	inputs := []run{
		{
			Input:  []byte{2, 16, 67, 97, 116, 101, 103, 111, 114, 121, 10, 98, 111, 111, 107, 115, 0},
			Output: map[string]string{"Category": "books"},
		},
		{
			Input:  []byte{4, 16, 67, 97, 116, 101, 103, 111, 114, 121, 10, 98, 111, 111, 107, 115, 8, 78, 97, 109, 101, 24, 72, 97, 114, 114, 121, 32, 112, 111, 116, 116, 101, 114, 0},
			Output: map[string]string{"Category": "books", "Name": "Harry potter"},
		},
	}

	for _, run := range inputs {
		bb := bytes.NewBuffer(run.Input)
		mp, err := ReadMapString(bb)
		if err != nil {
			t.Fatal(err)
		}

		for key, expected := range run.Output {
			val, has := mp[key]
			if !has {
				t.Fatalf("read output does not have the expected key: %s\n", key)
			}

			if val != expected {
				t.Fatalf("read output value does not match the expected output: %s, %s\n", expected, val)
			}
		}

		if bb.Len() != 0 {
			t.Fatal("not all bytes have been read from the byte buffer")
		}
	}
}

func TestWritingMapString(t *testing.T) {
	type run struct {
		Output []byte
		Input  map[string]string
	}

	inputs := []run{
		{
			Output: []byte{2, 16, 67, 97, 116, 101, 103, 111, 114, 121, 10, 98, 111, 111, 107, 115, 0},
			Input:  map[string]string{"Category": "books"},
		},
		{
			Output: []byte{4, 16, 67, 97, 116, 101, 103, 111, 114, 121, 10, 98, 111, 111, 107, 115, 8, 78, 97, 109, 101, 24, 72, 97, 114, 114, 121, 32, 112, 111, 116, 116, 101, 114, 0},
			Input:  map[string]string{"Category": "books", "Name": "Harry potter"},
		},
	}

	for _, run := range inputs {
		bb := bytes.NewBuffer(nil)
		err := WriteMapString(bb, run.Input)
		if err != nil {
			t.Fatal(err)
		}

		if bb.Len() != len(run.Output) {
			t.Fatalf("the returned byte buffer has an unexpected length: %b, %b\n", bb, run.Output)
		}

		for i, b := range bb.Bytes() {
			if b != run.Output[i] {
				t.Fatalf("unexpected byte encountered: %v, %v at index %d\n", b, run.Output[i], i)
			}
		}
	}
}

func BenchmarkReadingMapString(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteMapString(bb, map[string]string{"key": "val"})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadMapString(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingMapString(b *testing.B) {
	inputs := make([]map[string]string, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, map[string]string{"key": "val"})
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteMapString(bb, input)
	}
}
