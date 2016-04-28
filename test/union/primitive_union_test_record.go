package avro

import (
	"io"
)

type PrimitiveUnionTestRecord struct {
	UnionField UnionIntLongFloatDoubleStringBoolBytesNull
}

func DeserializePrimitiveUnionTestRecord(r io.Reader) (*PrimitiveUnionTestRecord, error) {
	return readPrimitiveUnionTestRecord(r)
}

func (r PrimitiveUnionTestRecord) Serialize(w io.Writer) error {
	return writePrimitiveUnionTestRecord(&r, w)
}
