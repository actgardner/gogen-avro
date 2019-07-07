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

	for expected, input := range inputs {
		bb := bytes.NewBuffer(input)

		result, err := ReadByte(bb)
		if err != nil {
			t.Fatal(err)
		}

		if string(result) != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %s recieved %s", input, expected, result)
		}

		if bb.Len() != 0 {
			t.Fatal("not all bytes have been read from the byte buffer")
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
				t.Fatalf("unexpected byte encountered: %v, %v at index %d\n", b, expected[i], i)
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

func TestReadingMapByte(t *testing.T) {
	type run struct {
		Input  []byte
		Output map[string][]byte
	}

	inputs := []run{
		{
			Input:  []byte{4, 16, 67, 97, 116, 101, 103, 111, 114, 121, 10, 98, 111, 111, 107, 115, 8, 67, 111, 100, 101, 24, 72, 97, 114, 114, 121, 32, 112, 111, 116, 116, 101, 114, 0},
			Output: map[string][]byte{"Category": []byte("books"), "Code": []byte("Harry potter")},
		},
		{
			Input:  []byte{2, 16, 67, 97, 116, 101, 103, 111, 114, 121, 10, 98, 111, 111, 107, 115, 0},
			Output: map[string][]byte{"Category": []byte("books")},
		},
	}

	for _, run := range inputs {
		bb := bytes.NewBuffer(run.Input)
		mp, err := ReadMapByte(bb)
		if err != nil {
			t.Fatal(err)
		}

		for key, expected := range run.Output {
			val, has := mp[key]
			if !has {
				t.Fatalf("read output does not have the expected key: %s\n", key)
			}

			for i, b := range val {
				if b != expected[i] {
					t.Fatalf("an unexpected byte was read at %d: %v, %v\n", i, expected[i], b)
				}
			}
		}

		if bb.Len() != 0 {
			t.Fatal("not all bytes have been read from the byte buffer")
		}
	}
}

func TestWritingMapByte(t *testing.T) {
	inputs := []map[string][]byte{
		{"Category": []byte("books"), "Code": []byte("Harry potter")},
		{"Category": []byte("books")},
	}

	for _, input := range inputs {
		bb := bytes.NewBuffer(nil)
		err := WriteMapByte(bb, input)
		if err != nil {
			t.Fatal(err)
		}

		mp, err := ReadMapByte(bb)
		if err != nil {
			t.Fatal(err)
		}

		if bb.Len() != 0 {
			t.Fatal("not all bytes have been read from the byte buffer")
		}

		for key, expected := range input {
			val, has := mp[key]
			if !has {
				t.Fatalf("an expected key is lost while writing %s\n", key)
			}

			for i, b := range val {
				if b != expected[i] {
					t.Fatalf("an unexpected byte was read at %d: %v, %v\n", i, expected[i], b)
				}
			}
		}
	}
}

func BenchmarkReadingMapByte(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteMapByte(bb, map[string][]byte{"key": []byte("val")})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadMapByte(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingMapByte(b *testing.B) {
	inputs := make([]map[string][]byte, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, map[string][]byte{"key": []byte("val")})
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteMapByte(bb, input)
	}
}
