package avro

import (
	"io"
)

type PrimitiveUnionTestRecord struct {
	UnionField UnionIntLongFloatDoubleStringBoolBytesNull
}

func (r PrimitiveUnionTestRecord) Serialize(w io.Writer) error {
	return writePrimitiveUnionTestRecord(&r, w)
}
