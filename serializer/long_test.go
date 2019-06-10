package serializer

import "testing"

func TestReadingLong(t *testing.T) {
	// expected result - input
	inputs := map[int64][]byte{
		20:                  []byte{40},
		30:                  []byte{60},
		9223372036854775807: []byte{254, 255, 255, 255, 255, 255, 255, 255, 255, 1}, // 64 bit signed int max
		60:                  []byte{120},
		15:                  []byte{30},
	}

	s := NewStream()
	i := NewLong(s)

	for expected, input := range inputs {
		go s.Write(input)

		result, err := i.Read()
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %d recieved %d", input, expected, result)
		}
	}
}
