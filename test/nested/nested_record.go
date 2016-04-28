package avro

import (
	"io"
)

type NestedRecord struct {
	StringField string
	BoolField   bool
	BytesField  []byte
}

func (r NestedRecord) Serialize(w io.Writer) error {
	return writeNestedRecord(&r, w)
}
