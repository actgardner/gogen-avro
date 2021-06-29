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
				t.Fatalf("unexpected byte encountered: %v, %v\n", b, expected[i])
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

	for expected, input := range inputs {
		bb := bytes.NewBuffer(input)

		result, err := ReadInt(bb)
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
				t.Fatalf("unexpected byte encountered: %v, %v at index %d\n", b, expected[i], i)
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

func TestReadingMapInt(t *testing.T) {
	type run struct {
		Input  []byte
		Output map[string]int32
	}

	inputs := []run{
		{
			Input:  []byte{4, 8, 67, 111, 100, 101, 200, 31, 16, 67, 97, 116, 101, 103, 111, 114, 121, 60, 0},
			Output: map[string]int32{"Category": 30, "Code": 2020},
		},
		{
			Input:  []byte{2, 16, 67, 97, 116, 101, 103, 111, 114, 121, 60, 0},
			Output: map[string]int32{"Category": 30},
		},
	}

	for _, run := range inputs {
		bb := bytes.NewBuffer(run.Input)
		mp, err := ReadMapInt(bb)
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

func TestWritingMapInt(t *testing.T) {
	inputs := []map[string]int32{
		{"Code": 2020, "Category": 30},
		{"Category": 30},
	}

	for _, input := range inputs {
		bb := bytes.NewBuffer(nil)
		err := WriteMapInt(bb, input)
		if err != nil {
			t.Fatal(err)
		}

		mp, err := ReadMapInt(bb)
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

func BenchmarkReadingMapInt(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteMapInt(bb, map[string]int32{"key": 100})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadMapInt(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingMapInt(b *testing.B) {
	inputs := make([]map[string]int32, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, map[string]int32{"key": 100})
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteMapInt(bb, input)
	}
}
