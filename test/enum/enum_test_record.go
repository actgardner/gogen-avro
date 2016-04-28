package avro

import (
	"io"
)

type EnumTestRecord struct {
	EnumField TestEnumType
}

func (r EnumTestRecord) Serialize(w io.Writer) error {
	return writeEnumTestRecord(&r, w)
}
