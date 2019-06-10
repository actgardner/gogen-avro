package serializer

import (
	"testing"
)

func TestReadingString(t *testing.T) {
	inputs := map[string][]byte{
		"john": []byte{8, 106, 111, 104, 110},
		"doe":  []byte{6, 100, 111, 101},
	}

	s := NewStream()
	b := NewString(s)

	for expected, input := range inputs {
		go s.Write(input)

		result, err := b.Read()
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %s recieved %s", input, expected, result)
		}
	}
}
