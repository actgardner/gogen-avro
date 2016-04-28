package avro

import (
	"io"
)

type ComplexUnionTestRecord struct {
	UnionField UnionNullArrayIntMapIntNestedUnionRecord
}

func (r ComplexUnionTestRecord) Serialize(w io.Writer) error {
	return writeComplexUnionTestRecord(&r, w)
}
