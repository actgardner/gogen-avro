package avro

import (
	"io"
)

type NestedTestRecord struct {
	NumberField *NumberRecord
	OtherField  *NestedRecord
}

func DeserializeNestedTestRecord(r io.Reader) (*NestedTestRecord, error) {
	return readNestedTestRecord(r)
}

func (r NestedTestRecord) Serialize(w io.Writer) error {
	return writeNestedTestRecord(&r, w)
}
