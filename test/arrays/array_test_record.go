package avro

import (
	"io"
)

type ArrayTestRecord struct {
	IntField    []int32
	LongField   []int64
	DoubleField []float64
	StringField []string
	FloatField  []float32
	BoolField   []bool
	BytesField  [][]byte
}

func (r ArrayTestRecord) Serialize(w io.Writer) error {
	return writeArrayTestRecord(&r, w)
}
