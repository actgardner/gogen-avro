package avro

import (
	"io"
)

type Parent struct {
	Children []*Child
}

func DeserializeParent(r io.Reader) (*Parent, error) {
	return readParent(r)
}

func (r Parent) Serialize(w io.Writer) error {
	return writeParent(&r, w)
}
