package avro

import (
	"io"
)

type NestedRecord struct {
	StringField string
	BoolField   bool
	BytesField  []byte
}

func DeserializeNestedRecord(r io.Reader) (*NestedRecord, error) {
	return readNestedRecord(r)
}

func (r NestedRecord) Serialize(w io.Writer) error {
	return writeNestedRecord(&r, w)
}
