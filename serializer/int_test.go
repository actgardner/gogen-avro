package serializer

import (
	"testing"
)

func TestReadingInt(t *testing.T) {
	// expected result - input
	inputs := map[int32][]byte{
		10:         []byte{20},
		30:         []byte{60},
		2147483647: []byte{254, 255, 255, 255, 15}, // 32 bit signed int max
		60:         []byte{120},
		15:         []byte{30},
	}

	s := NewStream()
	i := NewInt(s)

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
