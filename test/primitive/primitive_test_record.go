package avro

import (
	"io"
)

type PrimitiveTestRecord struct {
	IntField    int32
	LongField   int64
	FloatField  float32
	DoubleField float64
	StringField string
	BoolField   bool
	BytesField  []byte
}

func (r PrimitiveTestRecord) Serialize(w io.Writer) error {
	return writePrimitiveTestRecord(&r, w)
}
