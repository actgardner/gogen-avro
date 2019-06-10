package serializer

import "testing"

func TestHeaderMessageLength(t *testing.T) {
	inputs := map[int64][]byte{
		4:  []byte{8},
		1:  []byte{2},
		10: []byte{20},
	}

	s := NewStream()
	h := NewHeader(s)

	for expected, input := range inputs {
		go s.Write(input)

		result, err := h.ReadMessageLength()
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %d recieved %d", input, expected, result)
		}
	}
}
