package avro

import (
	"io"
)

type FixedTestRecord struct {
	FixedField TestFixedType
}

func (r FixedTestRecord) Serialize(w io.Writer) error {
	return writeFixedTestRecord(&r, w)
}
