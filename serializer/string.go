package serializer

import (
	"io"
)

// String ...
type String struct {
}

// Reader ...
func (s *String) Reader(r io.Reader) (empty []byte, err error) {
	return empty, err
}

// Writer ...
func (s *String) Writer(i []byte, w io.Writer) error {
	err := WriteLong(int64(len(i)), w)
	if err != nil {
		return err
	}

	_, err = w.Write(i)
	return err
}
