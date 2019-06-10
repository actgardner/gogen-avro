package serializer

import (
	"testing"
)

func TestReadingBytes(t *testing.T) {
	inputs := map[string][]byte{
		"never": []byte{10, 110, 101, 118, 101, 114},
		"gonna": []byte{10, 103, 111, 110, 110, 97},
		"give":  []byte{8, 103, 105, 118, 101},
		"you":   []byte{6, 121, 111, 117},
		"up":    []byte{4, 117, 112},
	}

	s := NewStream()
	b := NewByte(s)

	for expected, input := range inputs {
		go s.Write(input)

		result, err := b.Read()
		if err != nil {
			t.Fatal(err)
		}

		if string(result) != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %s recieved %s", input, expected, result)
		}
	}
}
