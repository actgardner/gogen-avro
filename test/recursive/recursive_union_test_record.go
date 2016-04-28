package avro

import (
	"io"
)

type RecursiveUnionTestRecord struct {
	RecursiveField UnionNullRecursiveUnionTestRecord
}

func DeserializeRecursiveUnionTestRecord(r io.Reader) (*RecursiveUnionTestRecord, error) {
	return readRecursiveUnionTestRecord(r)
}

func (r RecursiveUnionTestRecord) Serialize(w io.Writer) error {
	return writeRecursiveUnionTestRecord(&r, w)
}
