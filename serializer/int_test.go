package serializer

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestEncodingInt(t *testing.T) {
	// expected result - input
	inputs := map[uint64][]byte{
		10:         []byte{10},
		30:         []byte{30},
		2147483647: []byte{255, 255, 255, 255, 7}, // 32 bit signed int max
		60:         []byte{60},
		15:         []byte{15},
	}

	for input, expected := range inputs {
		bb := EncodeInt(10, input)

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

func TestReadingInt(t *testing.T) {
	// expected result - input
	inputs := map[int32][]byte{
		10:         []byte{20},
		30:         []byte{60},
		2147483647: []byte{254, 255, 255, 255, 15}, // 32 bit signed int max
		60:         []byte{120},
		15:         []byte{30},
	}

	r, w := io.Pipe()

	for expected, input := range inputs {
		go w.Write(input)

		result, err := ReadInt(r)
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %d recieved %d", input, expected, result)
		}
	}
}

func TestWritingInt(t *testing.T) {
	// expected result - input
	inputs := map[int32][]byte{
		10:         []byte{20},
		30:         []byte{60},
		2147483647: []byte{254, 255, 255, 255, 15}, // 32 bit signed int max
		60:         []byte{120},
		15:         []byte{30},
	}

	for input, expected := range inputs {
		r, w := io.Pipe()

		go func() {
			err := WriteInt(w, input)
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

func BenchmarkEncodingInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeInt(10, uint64(100))
	}
}

func BenchmarkReadingInt(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteInt(bb, 100)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadInt(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingInt(b *testing.B) {
	inputs := make([]int32, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, 100)
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteInt(bb, input)
	}
}
