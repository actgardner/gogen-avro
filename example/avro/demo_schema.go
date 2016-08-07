package avro

import (
	"io"
)

type DemoSchema struct {
	IntField    int32
	DoubleField float64
	StringField string
	BoolField   bool
	BytesField  []byte
}

func DeserializeDemoSchema(r io.Reader) (*DemoSchema, error) {
	return readDemoSchema(r)
}

func (r DemoSchema) Serialize(w io.Writer) error {
	return writeDemoSchema(&r, w)
}
