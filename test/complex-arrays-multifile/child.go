package avro

import (
	"io"
)

type Child struct {
	Name string
}

func DeserializeChild(r io.Reader) (*Child, error) {
	return readChild(r)
}

func (r Child) Serialize(w io.Writer) error {
	return writeChild(&r, w)
}
