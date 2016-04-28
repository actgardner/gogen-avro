package avro

import (
	"io"
)

type EnumTestRecord struct {
	EnumField TestEnumType
}

func DeserializeEnumTestRecord(r io.Reader) (*EnumTestRecord, error) {
	return readEnumTestRecord(r)
}

func (r EnumTestRecord) Serialize(w io.Writer) error {
	return writeEnumTestRecord(&r, w)
}
