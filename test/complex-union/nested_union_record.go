package avro

import (
	"io"
)

type NestedUnionRecord struct {
	IntField int32
}

func (r NestedUnionRecord) Serialize(w io.Writer) error {
	return writeNestedUnionRecord(&r, w)
}
