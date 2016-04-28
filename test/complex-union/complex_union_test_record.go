package avro

import (
	"io"
)

type ComplexUnionTestRecord struct {
	UnionField UnionNullArrayIntMapIntNestedUnionRecord
}

func DeserializeComplexUnionTestRecord(r io.Reader) (*ComplexUnionTestRecord, error) {
	return readComplexUnionTestRecord(r)
}

func (r ComplexUnionTestRecord) Serialize(w io.Writer) error {
	return writeComplexUnionTestRecord(&r, w)
}
