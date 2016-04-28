package avro

import (
	"io"
)

type NumberRecord struct {
	IntField    int32
	LongField   int64
	FloatField  float32
	DoubleField float64
}

func DeserializeNumberRecord(r io.Reader) (*NumberRecord, error) {
	return readNumberRecord(r)
}

func (r NumberRecord) Serialize(w io.Writer) error {
	return writeNumberRecord(&r, w)
}
