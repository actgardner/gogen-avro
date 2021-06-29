package serializer

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestReadingLong(t *testing.T) {
	// expected result - input
	inputs := map[int64][]byte{
		20:                  []byte{40},
		30:                  []byte{60},
		9223372036854775807: []byte{254, 255, 255, 255, 255, 255, 255, 255, 255, 1}, // 64 bit signed int max
		60:                  []byte{120},
		15:                  []byte{30},
	}

	for expected, input := range inputs {
		bb := bytes.NewBuffer(input)

		result, err := ReadLong(bb)
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

func TestWritingLong(t *testing.T) {
	// expected result - input
	inputs := map[int64][]byte{
		20:                  []byte{40},
		30:                  []byte{60},
		9223372036854775807: []byte{254, 255, 255, 255, 255, 255, 255, 255, 255, 1}, // 64 bit signed int max
		60:                  []byte{120},
		15:                  []byte{30},
	}

	for input, expected := range inputs {
		r, w := io.Pipe()

		go func() {
			err := WriteLong(w, input)
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

func BenchmarkReadingLong(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteLong(bb, 2147483648)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadLong(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingLong(b *testing.B) {
	inputs := make([]int64, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, 9223372036854775807)
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteLong(bb, input)
	}
}

func TestReadingMapLong(t *testing.T) {
	type run struct {
		Input  []byte
		Output map[string]int64
	}

	inputs := []run{
		{
			Input:  []byte{4, 8, 67, 111, 100, 101, 144, 3, 16, 67, 97, 116, 101, 103, 111, 114, 121, 200, 1, 0},
			Output: map[string]int64{"Category": 100, "Code": 200},
		},
	}

	for _, run := range inputs {
		bb := bytes.NewBuffer(run.Input)
		mp, err := ReadMapLong(bb)
		if err != nil {
			t.Fatal(err)
		}

		for key, expected := range run.Output {
			val, has := mp[key]
			if !has {
				t.Fatalf("read output does not have the expected key: %s\n", key)
			}

			if val != expected {
				t.Fatalf("read output value does not match the expected output: %d, %d\n", expected, val)
			}
		}

		if bb.Len() != 0 {
			t.Fatal("not all bytes have been read from the byte buffer")
		}
	}
}

func TestWritingMapLong(t *testing.T) {
	inputs := []map[string]int64{
		{"Category": 100, "Code": 200},
	}

	for _, input := range inputs {
		bb := bytes.NewBuffer(nil)
		err := WriteMapLong(bb, input)
		if err != nil {
			t.Fatal(err)
		}

		mp, err := ReadMapLong(bb)
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

			if val != expected {
				t.Fatalf("an unexpected value has been read at %s: %d, %d\n", key, expected, val)
			}
		}
	}
}

func BenchmarkReadingMapLong(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteMapLong(bb, map[string]int64{"key": 100})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadMapLong(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingMapLong(b *testing.B) {
	inputs := make([]map[string]int64, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, map[string]int64{"key": 100})
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteMapLong(bb, input)
	}
}
