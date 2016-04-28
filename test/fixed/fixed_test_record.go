package avro

import (
	"io"
)

type FixedTestRecord struct {
	FixedField TestFixedType
}

func DeserializeFixedTestRecord(r io.Reader) (*FixedTestRecord, error) {
	return readFixedTestRecord(r)
}

func (r FixedTestRecord) Serialize(w io.Writer) error {
	return writeFixedTestRecord(&r, w)
}
