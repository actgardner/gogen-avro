package avro

import (
	"io"
)

type NestedUnionRecord struct {
	IntField int32
}

func DeserializeNestedUnionRecord(r io.Reader) (*NestedUnionRecord, error) {
	return readNestedUnionRecord(r)
}

func (r NestedUnionRecord) Serialize(w io.Writer) error {
	return writeNestedUnionRecord(&r, w)
}
