package avro

import (
	"io"
)

type NestedTestRecord struct {
	NumberField *NumberRecord
	OtherField  *NestedRecord
}

func (r NestedTestRecord) Serialize(w io.Writer) error {
	return writeNestedTestRecord(&r, w)
}
