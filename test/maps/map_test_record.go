package avro

import (
	"io"
)

type MapTestRecord struct {
	IntField    map[string]int32
	LongField   map[string]int64
	DoubleField map[string]float64
	StringField map[string]string
	FloatField  map[string]float32
	BoolField   map[string]bool
	BytesField  map[string][]byte
}

func (r MapTestRecord) Serialize(w io.Writer) error {
	return writeMapTestRecord(&r, w)
}
