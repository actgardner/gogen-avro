package serializer

import (
	"errors"
	"testing"
)

func TestMessageLengthBytes(t *testing.T) {
	s := NewStream()
	b := NewByte(s)

	l := int64(64)

	go func() {
		err := b.WriteMessageLength(l)
		if err != nil {
			t.Fatal(err)
		}
	}()

	length, err := b.ReadMessageLength()
	if err != nil {
		t.Fatal(err)
	}

	if length != l {
		t.Fatal(errors.New("the expected message length is incorrect"))
	}
}

func TestReadingBytes(t *testing.T) {

}
