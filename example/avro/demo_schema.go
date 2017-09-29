/*
 * CODE GENERATED AUTOMATICALLY WITH gopkg.in/alanctgardner/gogen-avro.v5
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 *
 * SOURCE:
 *     example.avsc
 */

package avro

import (
	"gopkg.in/alanctgardner/gogen-avro.v5/container"
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

func NewDemoSchemaWriter(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &DemoSchema{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

func NewDemoSchema() *DemoSchema {
	v := &DemoSchema{}

	return v
}

func (r *DemoSchema) Schema() string {
	return "{\"fields\":[{\"name\":\"IntField\",\"type\":\"int\"},{\"name\":\"DoubleField\",\"type\":\"double\"},{\"name\":\"StringField\",\"type\":\"string\"},{\"name\":\"BoolField\",\"type\":\"boolean\"},{\"name\":\"BytesField\",\"type\":\"bytes\"}],\"name\":\"DemoSchema\",\"type\":\"record\"}"
}

func (r *DemoSchema) Serialize(w io.Writer) error {
	return writeDemoSchema(r, w)
}
